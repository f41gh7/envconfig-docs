# envconfig-docs

it provides doc printing feature for [envconfig package](https://github.com/kelseyhightower/envconfig)


## usage

 following struct field tags are supported:
 * default - defines default value for field
 * description - description of this field
 * envconfig - overrides generated envvar name for field
 * required - bool setting
 * split_words - splits struct field `SomeField` to `SOME_FIELD`, instead of `SOMEFIELD`.
 
 
Comments for struct field can be used for description if there is no description tag

install binary with command, you  have to install golang first.

```bash
go get -u github.com/f41gh7/envconfig-docs
```

create some file, let it be configs.go with struct, that configures your application:

```go
cat << 'EOF' > config.go
package conf

const prefixVar = "examplePref"

//genvars:true`
type BaseConfig struct {
	// Name of our main var
	// with multiline comment
	Name        string `default:"named"`
	// paths with array format
	MatchPaths []string `default:"path1,path2,path3"`
	MultiNested   struct {
		SomeField  string `default:"onelevelNest" description:"Nested field value with override" envconfig:"SIMPLE_FIELD"`
		InnerField struct {
			InnerCount     int    `default:"15" description:"inner field with int value"`
			TruncatedField string `default:"verylooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong value" description:"truncated, can be shown with flag -truncate=false"`
		}
	}
	// user password combination
	UserPassword map[string]string `default:"user:password1,user2:password2"`
}
EOF
```

Then generate docs for it:

```bash
envconfig-docs --input config.go
```
it will produce markdown output:

```markdown
# Auto Generated vars for package conf 
## updated at Mon Aug 17 00:50:40 UTC 2020 


| varible name | variable default value | variable required | variable description |
| --- | --- | --- | --- |
| examplePref_NAME | named | false | Name of our main varwith multiline comment |
| examplePref_MATCHPATHS | path1,path2,path3 | false | paths with array format |
| SIMPLE_FIELD | onelevelNest | false | Nested field value with override |
| examplePref__MULTINESTED_INNERFIELDINNERCOUNT | 15 | false | inner field with int value |
| examplePref__MULTINESTED_INNERFIELDTRUNCATEDFIELD | verylooooooooooooooooooooooooo | false | truncated, can be shown with flag -truncate=false |
| examplePref_USERPASSWORD | user:password1,user2:password2 | false | user password combination |

```

## configuration
 
 

 binary has following flags:

```bash
  -debugs string
        enables debug mode if not empty, debug will be written to stderr
  -input string
        input go file with config struct, by default conf.go
  -matchComment string
        struct comment line, that must be added to struct for match (default "//genvars:true")
  -output string
        out put file for variables
  -prefix string
        override prefix const value
  -prefixConst string
        prefix constant name at conf.go (default "prefixVar")
  -truncate
        truncate variable  value longer then 30 symbols (default true)
```