# OpenStreetMap Burmese Encoding

A QA tool for the [OpenStreetMap](http://openstreetmap.org) database to detect the non-Unicode [Zawgyi text encoding](https://en.wikipedia.org/wiki/Zawgyi_font) for Burmese (Myanmar) language.

Because Zawgyi occupies the same codepoints in Unicode, it must be classified using probability.

This tool uses:

* Google's [myanmar-tools](https://github.com/google/myanmar-tools) classifier
* [paulmach/osm](https://github.com/paulmach/osm) for parsing OSM
* [Rabbit](github.com/Rabbit-Converter/Rabbit-Go) for suggesting Zawgyi->Unicode conversions
