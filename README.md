# OpenStreetMap Burmese Encoding

A QA tool for the [OpenStreetMap](http://openstreetmap.org) database to detect the non-Unicode [Zawgyi text encoding](https://en.wikipedia.org/wiki/Zawgyi_font) for Burmese (Myanmar) language.

View the daily-generated data in the Viewer at [bdon.github.io/OpenStreetMap-BurmeseEncoding](https://bdon.github.io/OpenStreetMap-BurmeseEncoding)

Without a consistent text encoding, it is impossible for application developers to display text correctly across the world. Note that software like Mapnik and Canvas-based renderers can draw Burmese Unicode text correctly; [MapLibre GL JS](https://github.com/wipfli/about-text-rendering-in-maplibre) currently cannot.

Because Zawgyi occupies the same codepoints in Unicode, it must be classified using probability.

This tool uses:

* Google's [myanmar-tools](https://github.com/google/myanmar-tools) classifier
* [paulmach/osm](https://github.com/paulmach/osm) for parsing OSM
* [Rabbit](https://github.com/Rabbit-Converter/Rabbit-Go) for suggesting Zawgyi->Unicode conversions
  
## How to Contribute

Do not automate edits to OpenStreetMap data.

1. Create an account on [openstreetmap.org](http://openstreetmap.org)
2. Click the link to an OSM Object on the [Viewer](https://bdon.github.io/OpenStreetMap-BurmeseEncoding)
3. Edit text using the web editor
4. Add `#zawgyi` to your changeset description.
   
## License

BSD-2-Clause, except for [Rabbit](https://github.com/Rabbit-Converter/Rabbit-Go) which is pending license.
