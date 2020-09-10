# tymlate
tymlate is a folder structure template engine based on go template

## Usage

```text
Usage:
  tymlate [flags]

Flags:
  -c, --configuration string   path to configuration file
  -d, --datasource strings     Datasource in name=file format
  -h, --help                   help for tymlate
  -s, --source string          path to template source folder
  -t, --target string          path to target folder

```

In order to use `tymlate` you need:
- folder containing source tree with files to be templated
- target folder where the files will be generated
- configuration file inside of source folder with name 
`.tymlate.yml` or provided with `-c` flag
- optionally data yml files 

## Configuration file

Here is an example configuration file:
```yaml
data: #direct data
  meta:
    package: "github.com/lorands/tymlate"
    name: MyEzample
  Dat1: # the name of the map
    key1: "value1"
    key2: 22
    key3: true
  Dat2:
    key11: "value eleven"
    key2: 23
    key3: false
include:
  cfgName1: subcfg1.yml
  cfgNameTwo: subconf/subcfg2.yml
templates:
  includes:
  - ".*\\.tpl"
  - ".*\\.go"
  - LICENSE
  excludes:
  - ".*\\.inv"
  suffix: ".tpl"
```
Where:

1. `data` (optional): contains the variables to be used in context,
from the example you would access {{.meta.name}} from your template
1. `include` (optional): include additional yaml data files,
the key will be the key from the yaml file (e.g. `cfgName1`)
1. `templates` (required)
    1. `includes` - files to include (regexp)
    1. `excludes` - files to ignore (regexp)
    1. `suffix` (required): the template file suffix


## Step by step usage

1. Create a template folder structure. 
1. Create a configuration that defines how to process
 the structure.
1. Define folder structure
1. Define context data file(s)
1. Run tymlate with configuration

Let's assume you would like to generate this structure 
each time you have a new project:

```
.
|___ cmd /
|  |___ main.go
|___ model /
|  |___ model.go
|___ README.md
|___ LICENSE
|___ Dockerfile
|___ go.mod

```

### Template folder structure

Just create a structure that will be mirrored to target.

Decide on what will be the template file extension, 
we use `.tpl` in this example (see: `tymlate/generator/testdata/source`).

```
.
|___ cmd /
|  |___ main.go.tpl
|___ model /
|  |___ model.go.tpl
|___ README.md.tpl
|___ LICENSE
|___ Dockerfile.tpl
|___ go.mod.tpl
```

Each template file (.tpl in this case), is a go template file.

## It's just a go template

Each template file is "just" a go template.

## Template functions included

We included Sprig (http://masterminds.github.io/sprig/)
functions for your convenience. 

## Command line (cli) examples

The simplest case:
```shell script
tymlate -s path/to/source -t path/to/target
```
In this case the configuration should be provided
in `.tymlate.yml` file inside the source.

Simple case, with an external config: 
```shell script
tymlate -c path/to/config.yml -s path/to/source -t path/to/target
```

Additional data sources
```shell script
tymlate -c path/to/config.yml -d ds=path/to/data.yml -s path/to/source -t path/to/target
```

More data sources
```shell script
tymlate -c path/to/config.yml -d ds1=path/to/data1.yml -d ds2=path/to/data2.yml -s path/to/source -t path/to/target
```

