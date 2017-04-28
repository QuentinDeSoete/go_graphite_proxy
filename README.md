# go_graphite_proxy v0.1.0 [![Build Status](https://travis-ci.org/QuentinDeSoete/go_graphite_proxy.svg?branch=master)](https://travis-ci.org/QuentinDeSoete/go_graphite_proxy)

This is a proxy which comes between [Grafana](https://grafana.com) and [Graphite](https://graphiteapp.org).
It will authentificate the users from grafana on graphite like that your [Grafana](https://grafana.com) users will not be able to pull all [Graphite](https://graphiteapp.org) metrics.

## Usage

Just run the binaire as a daemon and set up your datasource like bellow in [Grafana](https://grafana.com):
<p><img src="img/datasource.png"></p>

## Contributing

Feel free to make a pull request.

## License

```
MIT License

Copyright (c) 2017 Quentin De Soete

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
