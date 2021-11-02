# RPM Package

To build an rpm package for CAS:

```sh
rpmdev-setuptree
cp cas.spec ~/rpmbuild/SPECS
cp man.patch ~/rpmbuild/SOURCES
wget https://github.com/codenotary/cas/archive/refs/tags/v1.0.0.tar.gz -O ../SOURCES/v1.0.0.tar.gz
cd ~/rpmbuild/SPECS
rpmbuild -ba cas.spec
```

### Note
File `man.patch` contains the patch needed to generate man pages (not included in
1.0.0).
