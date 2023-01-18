[![CI](https://github.com/bitnami/render-template/actions/workflows/main.yml/badge.svg)](https://github.com/bitnami/render-template/actions/workflows/main.yml)

# render-template

This tools allows rendering Handlebars 3.0 templates, using as context data the current environment variables or a provided data file.

# Basic usage

```console
$ render-template --help
Usage:
  render-template [OPTIONS] [template-file]

Application Options:
  -f, --data-file=DATA_FILE    Properties file containing the replacements for the template

Help Options:
  -h, --help                   Show this help message

Arguments:
  template-file:               File containing the template to render. Its contents can be also passed through stdin
```

The tool supports rendering templates from a file or from stdin (for convenience).

The source data is taken from the environment variables or from a data file, with properties-file format (key=value, a line for each pair). When a variable is defined both as an environment variable and in the data file, the latter will take precedence.

# Examples

## Render data from template file with environment variables

```console
# Create the template
$ echo 'hello {{who}}' > template.tpl
# Render it without 'who' variable
$ render-template template.tpl
hello 
# Render it with 'who' variable defined
$ who=bitnami render-template template.tpl
hello bitnami
```

## Render data from stdin with environment variables

```console
$ log_file=/tmp/stout.log port=8080 pid_file=/tmp/my.pid render-template <<"EOF"
# My servide log file
log_file "{{log_file}}"

# HTTP port
port {{port}}

# My service pid file
pid_file "{{pid_file}}"
EOF
```

Outputs:

```
# My servide log file
log_file "/tmp/stout.log"

# HTTP port
port 8080

# My service pid file
pid_file "/tmp/my.pid"
```

## Render data from stdin with data file

```console
# write data file
$ cat > data.properties <<"EOF"
log_file=/tmp/stout.log
port=8080
pid_file=/tmp/my.pid 
EOF

$ render-template --data-file ./data.properties <<"EOF"
# My servide log file
log_file "{{log_file}}"

# HTTP port
port {{port}}

# My service pid file
pid_file "{{pid_file}}"
EOF
```

Outputs:

```
# My servide log file
log_file "/tmp/stout.log"

# HTTP port
port 8080

# My service pid file
pid_file "/tmp/my.pid"
```

## Overriding environment variables in data file

```console
# Lets define some environment variables
$ export name=foo
$ export company=bar
$ export year=3000

# And write a template
$ cat > data.tpl <<"EOF"
{{name}} works at {{company}}
since {{year}}
EOF

# Rendering from the environment would yield
$ render-template data.tpl
foo works at bar
since 3000

# But we can override it from a data file, either partially, to get a mix:

$ echo "name=mike" > data.properties
$ render-template --data-file data.properties data.tpl
mike works at bar
since 3000

# Or completely:

$ cat > data.properties <<"EOF"
name=mike
company=Bitnami
year=2010
EOF

$ render-template --data-file data.properties data.tpl
mike works at Bitnami
since 2010
```

## Using helpers

The tool supports all the standard handlebars helpers: https://handlebarsjs.com/builtin_helpers.html

```console
$ render-template <<"EOF" 
 {{#if author}}
{{firstName}} {{lastName}}
 {{else}}
Unknown Author
 {{/if}}
EOF

# Which outputs
Unknown Author

$ author=me firsName=foo lastName=bar render-template <<"EOF" 
 {{#if author}}
{{firstName}} {{lastName}}
 {{else}}
Unknown Author
 {{/if}}
EOF

# Outputs:
foo bar
```

In addition, it includes a few custom helpers:

### json_escape

The json_escape helper converts the  provided value into a valid JSON string
```console
$ export VALUE='this is "a string", with quoting

and some line breaks'
```
Without the helper:

```console
$ render-template <<<'VALUE={{VALUE}}'
VALUE=this is "a string", with quoting

and some line breaks
```

Using the helper:
```console
$ render-template <<<'VALUE={{json_escape VALUE}}'
VALUE="this is \"a string\", with quoting\n\nand some line breaks"
```

### quote

The quote helper Quotes a string

Without the helper:
```console
$ ARG1="some arg" ARG2="some other \"arg\"" render-template <<"EOF"
ARG1={{ARG1}} ARG2={{ARG2}}
EOF
ARG1=some arg ARG2=some other "arg"
```

With the helper

```console
ARG1="some arg" ARG2="some other \"arg\"" render-template <<"EOF"
ARG1={{quote ARG1}} ARG2={{quote ARG2}}
EOF
ARG1="some arg" ARG2="some other \"arg\""
```

### or

This helper allows using the "or" logical operation over two values (a value will be true if not empty)

To render a block when either "firstName" or "lastName" values ar not empty:

```console
$ cat > data.tpl <<"EOF" 
{{#if (or firstName lastName)}}
{{firstName}} {{lastName}}
{{else}}
Unknown Author
{{/if}}
EOF

$ render-template data.tpl
Unknown Author

$ firstName=foo render-template data.tpl
foo

$ lastName=bar render-template data.tpl
bar
 
$ firstName=foo lastName=bar render-template data.tpl
foo bar
```

This helper can also be used to provide defaults for your template variables:

```console
# No value provided, so we fallback to the second "or" argument
$ render-template <<<'VALUE={{or ENV_VALUE "default value"}}'
VALUE=default value

# ENV_VALUE is defined, so we take it
$ ENV_VALUE="customized value" render-template <<<'VALUE={{or ENV_VALUE "default value"}}'
VALUE=customized value
```

## License

Copyright &copy; 2023 Bitnami

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.

You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and limitations under the License.
