/*
 * Copyright (c) 2018-2021 Codenotary, Inc. All Rights Reserved.
 * This software is released under Apache License 2.0.
 * The full license information can be found under:
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 */

package executor

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/pkg/stdcopy"
)

// timeout for kill ing the container if it doesn't properly shuts down
var timeout = 5 * time.Second

type DockerExecutor struct {
	ctx    context.Context
	client *docker.Client
	contID string
}

// NewDockerExecutor starts a new container and returns an executor for the container
func NewDockerExecutor(image string) (Executor, error) {
	ctx := context.Background()

	dockerClient, _ := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	cont, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image:       image,
			Entrypoint:  []string{"/bin/sh"},
			AttachStdin: true,
			Tty:         true,
			OpenStdin:   true,
		},
		nil, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("cannot create container: %w", err)
	}

	err = dockerClient.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{})
	if err != nil {
		if errdefs.IsInvalidParameter(err) {
			return nil, fmt.Errorf("cannot use container without shell. Images 'from scratch' are not supported")
		}
		return nil, err
	}

	return DockerExecutor{
		ctx:    ctx,
		client: dockerClient,
		contID: cont.ID,
	}, nil
}

// Exec executes a command inside the container
func (e DockerExecutor) Exec(cmd []string) ([]byte, []byte, int, error) {
	var stdOut, stdErr bytes.Buffer
	exec, err := e.client.ContainerExecCreate(e.ctx, e.contID, types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return nil, nil, 0, err
	}

	hijacked, err := e.client.ContainerExecAttach(e.ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, nil, 0, err
	}
	defer hijacked.Close()

	done := make(chan error)

	go func() {
		_, err := stdcopy.StdCopy(&stdOut, &stdErr, hijacked.Reader)
		done <- err
	}()

	select {
	case err = <-done:
		if err != nil {
			return nil, nil, 0, err
		}
		break
	case <-e.ctx.Done():
		return nil, nil, 0, e.ctx.Err()
	}

	// check exit code
	res, err := e.client.ContainerExecInspect(e.ctx, exec.ID)
	if err != nil {
		return nil, nil, 0, err
	}

	return stdOut.Bytes(), stdErr.Bytes(), res.ExitCode, nil
}

// Close stops previously started container by sending "exit" to shell - it is much faster than
// stopping with container API
func (e DockerExecutor) Close() error {
	hijacked, err := e.client.ContainerAttach(e.ctx, e.contID, types.ContainerAttachOptions{
		Stdin:  true,
		Stream: true,
	})
	if err != nil {
		// force stop
		fmt.Printf("cannot attach to container: %s", err)
		return e.client.ContainerStop(e.ctx, e.contID, &timeout)
	}
	_, err = hijacked.Conn.Write([]byte("exit\n"))
	if err != nil {
		// force stop
		fmt.Printf("cannot send command to container: %s", err)
		return e.client.ContainerStop(e.ctx, e.contID, &timeout)
	}
	hijacked.Conn.Close()
	doneCh, errCh := e.client.ContainerWait(e.ctx, e.contID, container.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		return err
	case <-doneCh:
		return nil
	}
}

// Read reads the file from container and returns its content
func (e DockerExecutor) ReadFile(path string) ([]byte, error) {
	reader, _, err := e.client.CopyFromContainer(e.ctx, e.contID, path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	tr := tar.NewReader(reader)
	_, err = tr.Next()
	if err != nil {
		return nil, fmt.Errorf("error processing file from container: %w", err)
	}

	return ioutil.ReadAll(tr)
}

// ReadDir reads the files from container directory and returns the content as TAR stream
func (e DockerExecutor) ReadDir(path string) (io.ReadCloser, error) {
	reader, _, err := e.client.CopyFromContainer(e.ctx, e.contID, path)
	return reader, err
}
