package docker

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	rpmdb "github.com/anchore/go-rpmdb/pkg"

	"github.com/codenotary/cas/pkg/bom/artifact"
	"github.com/codenotary/cas/pkg/bom/executor"
)

type rpm struct {
	db     *rpmdb.RpmDB
	byFile map[string]*rpmdb.PackageInfo
}

var hashTypeMaps = map[rpmdb.DigestAlgorithm]artifact.HashType{
	rpmdb.PGPHASHALGO_MD5:    artifact.HashMD5,
	rpmdb.PGPHASHALGO_SHA1:   artifact.HashSHA1,
	rpmdb.PGPHASHALGO_MD2:    artifact.HashMD2,
	rpmdb.PGPHASHALGO_SHA256: artifact.HashSHA256,
	rpmdb.PGPHASHALGO_SHA384: artifact.HashSHA384,
	rpmdb.PGPHASHALGO_SHA512: artifact.HashSHA512,
	rpmdb.PGPHASHALGO_SHA224: artifact.HashSHA224,
}

func (pkg rpm) Type() string {
	return RPM
}

func (pkg *rpm) AllPackages(e executor.Executor, output artifact.OutputOptions) ([]artifact.Dependency, error) {
	var err error
	if pkg.db == nil {
		pkg.db, err = openDb(e)
		if err != nil {
			return nil, err
		}
	}

	pkgList, err := pkg.db.ListPackages()
	if err != nil {
		return nil, fmt.Errorf("cannot read RPM database: %w", err)
	}

	res := make([]artifact.Dependency, 0, len(pkgList))
	for _, p := range pkgList {
		hashtype, ok := hashTypeMaps[p.DigestAlgorithm]
		if !ok {
			hashtype = artifact.HashInvalid
		}
		hash, err := combineHashesFromSlice(p.Files, p.Name)
		if err != nil {
			fmt.Printf("Cannot combine hashes: %v\n", err)
			// ignore error
		}
		if hash == "" {
			continue
		}
		if p.License == "" {
			p.License = "NONE"
		}
		res = append(res, artifact.Dependency{
			Name:     p.Name,
			Version:  p.Version + "-" + p.Release,
			HashType: hashtype,
			Hash:     hash,
			License:  p.License,
			Type:     artifact.DepDirect,
		})
	}

	return res, nil
}

func openDb(e executor.Executor) (*rpmdb.RpmDB, error) {
	buf, err := e.ReadFile("/var/lib/rpm/Packages")
	if err != nil {
		return nil, fmt.Errorf("error reading file from container: %w", err)
	}

	f, err := ioutil.TempFile("", "cas.rpmdb")
	if err != nil {
		return nil, fmt.Errorf("cannot create temporary file: %w", err)
	}
	rpmFile := f.Name()
	defer os.Remove(rpmFile) // clean up

	_, err = f.Write(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot write temporary file: %w", err)
	}
	f.Close()

	db, err := rpmdb.Open(rpmFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read RPM database: %w", err)
	}

	return db, err
}

func combineHashesFromSlice(files []rpmdb.FileInfo, pkgName string) (string, error) {
	var hash []byte
	for _, file := range files {
		if file.Digest == "" {
			// some files don't have digests, like symbolic links
			continue
		}
		comp, err := hex.DecodeString(file.Digest)
		if err != nil {
			return "", fmt.Errorf("malformed hash for package %s", pkgName)
		}
		if hash == nil {
			hash = comp
		} else {
			if len(comp) != len(hash) {
				// should never happen - all hashes must be of the same length
				return "", fmt.Errorf("malformed hash for package %s", pkgName)
			}
			// XOR hash
			for i := 0; i < len(hash); i++ {
				hash[i] ^= comp[i]
			}
		}
	}

	return hex.EncodeToString(hash), nil
}
