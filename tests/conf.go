package conf

import "time"

const prefixVar = "examplePref"

//BaseConfig - example
//genvars:true`
type BaseConfig struct {
	// Name of our main var
	// with multiline comment
	Name        string `default:"named"`
	NestedField struct {
		// TODO must be ignored
		// must present
		Listing string `default:"all available"`
		Enabled bool   `default:"false" description:"enables something" required:"true"`
	}
	AnotherOption string `default:"AnotherValue" description:"some other value"`
	// its default timeout
	// for our oper
	Timeout time.Duration `default:"10s"`
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