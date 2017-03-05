# BandiKoot

Kolla arch is nice but its got a couple shortcomings
- resulting images are huge
- layering is a disaster
- too much build cruft is left lying around (i.e. toolchain)

Lets fix it with Alpine!!!

# Desired Goals
- need some kind of templating / preprocessing

- static uid / gid across the board
- sudo user perm setup
- build in unit tests

- users should be able to easily override default yaml configs
- override pip versions
- add / configure additional packages
- start / end blocks for package

- CI friendly

- everything is HA capable

# Implemented
- all openstack based on python 2.7
- pip driven source installs only
- 1 layer, always cleanup build deps
- need CMD binary. (not really, assume helm will handle all that)
- Makefile
- single config file, multiple supported also (env.sh)
- upstream alpine images whenever possible (3.5 across the board?)
- select git repo
- select branch
