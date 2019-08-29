# webshorten
n URL shortening web service using the only standard library and clearuri.com API
=====

`webshorten` is the tool to shorten URL by sending POST reqs to a specified URL.

Installation
------------

Install `webshorten` from sources, by running:

```sh
git clone https://github.com/PhilLar/webshorten.git
cd webshorten
go install ./cmd/webshorten
```

Usage
-----
You can shorten the URL:
```sh
webshorten
```
OR
accompany with a flag
``sh
webshorten -port 5001
```
OR
accompany with PORT env var specified
``sh
PORT=8080 webshorten
```
(default port value is ":5000")

Contribute
----------
- Issue Tracker: https://github.com/PhilLar/webshorten/issues
- Source Code: https://github.com/PhilLar/webshorten

License
--------
[WTFPL 2.0](https://wtfpl2.com/)
