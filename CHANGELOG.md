# Changelog

## 1.0.0 (2026-09-20)


### Features

* add converter package (options, scanner, dnglab wrapper) ([18db956](https://github.com/Strobotti/dnglab-gui/commit/18db9567e96f20f715a7114233e726feb6550873))
* add main window UI ([9e4b63d](https://github.com/Strobotti/dnglab-gui/commit/9e4b63d70ddc66ad0e4c480849914f2b2e6efb3c))
* add progress dialog and about dialog ([416d26f](https://github.com/Strobotti/dnglab-gui/commit/416d26f6b0e7593da95ee43a424b338f23699ee8))
* integrate application icon using fyne bundle ([668b80f](https://github.com/Strobotti/dnglab-gui/commit/668b80f01dcdf747f095e15fdba1b4d43a53d3b8))
* wire up conversion logic end to end ([70837d0](https://github.com/Strobotti/dnglab-gui/commit/70837d0e9311fb4ba3d15b15c037db7e05e0e625))


### Bug Fixes

* address review findings - batch conversion, fyne.Do safety, settings persistence, crop flag ([8298285](https://github.com/Strobotti/dnglab-gui/commit/82982855761818536d374bc1e7e8e358f0f29a77))
* Convert test, progress indicator, goroutine-safety comment scope ([b534bf3](https://github.com/Strobotti/dnglab-gui/commit/b534bf31026987e29e84eb1a3a50e21ed41f47fe))
* improve UI section labels and layout ([bf932f5](https://github.com/Strobotti/dnglab-gui/commit/bf932f568d05dadbcfc0f7ca03cb300b698e56f8))
* increase main window height to avoid scrollbar ([7ce7100](https://github.com/Strobotti/dnglab-gui/commit/7ce7100360960d077ed5a7f995da51cd03722a21))
* make output folder dropdown expand to fill available width ([5330a01](https://github.com/Strobotti/dnglab-gui/commit/5330a013d44543761eef177570880f287ea59d89))
* persist recursive setting, eliminate double ScanForRawFiles ([296eebe](https://github.com/Strobotti/dnglab-gui/commit/296eebee34a8c6e7f015bfba2b1bab0238cab8af))
* put compression and crop labels inline with their radio groups ([f3a9f19](https://github.com/Strobotti/dnglab-gui/commit/f3a9f198eafdfd56985c27965a1fa3aa95f15eb1))
* use app.NewWithID() to enable Fyne Preferences API ([8099882](https://github.com/Strobotti/dnglab-gui/commit/809988264a43d08bb5e66b20c7eca87442c4112f))
