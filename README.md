### Q?

Show me the mean durations of successful requests to /events broken down by app id


### A?

Given some ruleset:

```
{
    "rule": {
        "type": "and",
        "rules": [
            { "type": "all_match" },
            { "type": "field_match", "name": "status", "value": 200 },
            { "type": "field_match", "name": "path", "value": "/mypath" }
        ]
    },
    "group": "appname",
    "calculation": {
        "type": "mean",
        "name": "duration"
    }
}
```

And some logs:

```
$ pingu < resources/full.log -l 10
Grouping by appname:

myapp2: 112.25
myapp: 108.09
myapp3: 106.43
myapp4: 100.65
myapp7: 93.69
myapp10: 93.55
myapp9: 86.47
myapp6: 85.2
myapp8: 81.485
myapp5: 79.95
```
