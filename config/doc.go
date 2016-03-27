/*
The config package contains the API needed to make your plugin configurable.

Making your plugin available in the config file

To be used in the config file, your plugin must implement the Builder
interface, where it will be provided a Config object, which allows the plugin
programmer to load JSON into an annotated struct of their choice, using
ConfigSet's ParseConfig() method. From there, you can configure your
Producer as you see fit with your user-provided configuration.

I have followed an API similar to that of golang's own database/sql.
When writing your plugin, you will need to register yourself in an init()
method in your package, and then the ConfigSet object will be aware of your
plugin and will use it when given the corresponding key.

After registering, you need to include your package as an anonymous import in
the main function which parses the config file - you can find it in cmd/goi3bar

Builder

Build() is the method that is required to implement the Builder interface.
It is provided a Config object, which offers ParseConfig() to cast the options
subtree of your plugin's JSON block to a struct a-la json.Unmarshal. The Config
object also has the raw interface{} of the options subtree, allowing you to
inspect it manually, so you can provide a flexible API without the overhead of
remarshalling through trial-and-error.
See the network builders for an example of this.

I have defined an interface instead of a simple function typedef because there
are some opportunities for code re-use. See the CPU builder for an example.

Get involved

If you write something you think is really nifty, you are of course welcome
to submit a pull request to me for consideration to have it included in the
default distribution.
*/
package config
