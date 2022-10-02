# go-randomize-workinghour

## build

```
% go build
```

## settings

create config.yml for setting both start and end time of working hour.  
this tool generates random minute in 20 for both times.

```yaml
from:      # time to start your work.
  hour: 9
  min: 40
to:        # time to end your work
  hour: 19
  min: 0
```

## how to use

```bash
# run at 2022-10-02
% ./go-randomize-workinghour
Generated working times in 2022/10, lets copy into Google Spreadsheet

% ./go-randomize-workinghour 2022-09
Generated working times in 2022/09, lets copy into Google Spreadsheet
```