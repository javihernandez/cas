Name:           cas
Version:        1.0.0
Release:        1%{?dist}
Summary:        Community Attestation Service

License:        ASL 2.0
URL:            https://github.com/codenotary/cas
Source0:        https://github.com/codenotary/cas/archive/refs/tags/v%{version}.tar.gz
Patch0:         man.patch

BuildRequires:  make golang

%description
Give any digital asset a meaningful, globally-unique, immutable identity 
that is authentic, verifiable, traceable from anywhere.

# workaround for missing GCC builid
%global _missing_build_ids_terminate_build 0
%global debug_package %{nil}

%prep
%setup -q
%patch0 -p1

%build
make
make docs/cmd

%install
mkdir -p %{buildroot}%{_bindir} %{buildroot}%{_mandir}/man1
install -p -m 755 %{name} %{buildroot}%{_bindir}
install -p -m 644 docs/man/*.1 %{buildroot}%{_mandir}/man1

%clean
rm -rf %{buildroot}

%files
%{_bindir}/cas
%{_mandir}/man1/*.1.gz

%doc README.md CONTRIBUTING.md CHANGELOG.md
%license LICENSE

%changelog
* Fri Oct 29 2021 simone 1.0.0-1
- Initial version of the package
