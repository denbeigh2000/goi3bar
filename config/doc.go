/*
	The config package contains the API needed to make your plugin configurable.

	To be used in the Config file, your plugin must implement the Builder
	interface, where it will be provided a Config object, which allows the plugin
	programmer to load JSON into an annotated struct of their choice, using
	ConfigSet's ParseConfig() method. From there, you can configure your
	Producer as you see fit with your user-provided configuration.

	I have followed an API similar to that of golang's own database/sql.
	When writing your plugin, you will need to register yourself in an init()
	method in your package, and then the ConfigSet object will be aware of your
	plugin and will use it when given the corresponding key.
	This has the pleasant benefit that the user can simply add your custom package
	to the config parser mainfile as an anonymous import (as you would a sql
	driver), compile, and provide a config file.

	If you write something you think is really nifty, you are of course welcome
	to submit a pull request to me at the home repository and have it included in
	the default distribution.
*/
package config
