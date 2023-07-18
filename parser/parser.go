package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type CronField struct {
	Name     string
	MinValue int
	MaxValue int
	Value    []int
}

type CronExpression struct {
	Second     CronField
	Minute     CronField
	Hour       CronField
	DayOfMonth CronField
	Month      CronField
	DayOfWeek  CronField
}

const (
	MON = 1
	TUE = 2
	WED = 3
	THU = 4
	FRI = 5
	SAT = 6
	SUN = 0
)

func Trans(str string) (Int int, err error) {
	switch str {
	case "MON":
		return MON, nil
	case "TUE":
		return TUE, nil
	case "WED":
		return WED, nil
	case "THU":
		return THU, nil
	case "FRI":
		return FRI, nil
	case "SAT":
		return SAT, nil
	case "SUN":
		return SUN, nil
	default:
		result, err := strconv.Atoi(str)
		return result, err
	}
}

func NewCronExpression() *CronExpression {
	return &CronExpression{
		Second: CronField{
			Name:     "second",
			MinValue: 0,
			MaxValue: 59,
		},
		Minute: CronField{
			Name:     "minute",
			MinValue: 0,
			MaxValue: 59,
		},
		Hour: CronField{
			Name:     "hour",
			MinValue: 0,
			MaxValue: 23,
		},
		DayOfMonth: CronField{
			Name:     "day of month",
			MinValue: 1,
			MaxValue: 31,
		},
		Month: CronField{
			Name:     "month",
			MinValue: 1,
			MaxValue: 12,
		},
		DayOfWeek: CronField{
			Name:     "day of week",
			MinValue: 0,
			MaxValue: 6,
		},
	}
}

func (ce *CronExpression) Parse(expression string) error {
	fields := strings.Fields(expression)
	if len(fields) != 6 {
		return fmt.Errorf("invalid cron expression")
	}

	err := ce.parseField(fields[0], &ce.Second)
	if err != nil {
		return err
	}

	err = ce.parseField(fields[1], &ce.Minute)
	if err != nil {
		return err
	}

	err = ce.parseField(fields[2], &ce.Hour)
	if err != nil {
		return err
	}

	if fields[3] == "?" && fields[5] == "?" {
		return fmt.Errorf("invalid cron expression")
	}

	if fields[3] != "?" {
		err = ce.parseField(fields[3], &ce.DayOfMonth)
		if err != nil {
			return err
		}
	}

	err = ce.parseField(fields[4], &ce.Month)
	if err != nil {
		return err
	}

	if fields[5] != "?" {
		err = ce.parseField(fields[5], &ce.DayOfWeek)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ce *CronExpression) parseField(field string, cf *CronField) error {
	if field == "*" || field == "?" {
		for i := cf.MinValue; i <= cf.MaxValue; i++ {
			cf.Value = append(cf.Value, i)
		}
		return nil
	}

	parts := strings.Split(field, ",")
	for _, part := range parts {
		stepParts := strings.Split(part, "/")

		// .../...
		if len(stepParts) == 2 {
			step, err := Trans(stepParts[1])
			if err != nil || step < 1 || step > cf.MaxValue-cf.MinValue+1 {
				return fmt.Errorf("invalid %s field", cf.Name)
			}

			rangeParts := strings.Split(stepParts[0], "-")

			// ...-.../...
			if len(rangeParts) == 2 {
				start, err := Trans(rangeParts[0])
				if err != nil || start < cf.MinValue || start > cf.MaxValue {
					return fmt.Errorf("invalid %s field", cf.Name)
				}

				end, err := Trans(rangeParts[1])
				if err != nil || end < cf.MinValue || end > cf.MaxValue || end < start {
					return fmt.Errorf("invalid %s field", cf.Name)
				}

				for i := start; i <= end; i += step {
					cf.Value = append(cf.Value, i)
				}

			} else if stepParts[0] == "*" {
				for i := cf.MinValue; i <= cf.MaxValue; i += step {
					cf.Value = append(cf.Value, i)
				}
			} else {
				return fmt.Errorf("invalid %s field", cf.Name)
			}
		} else if len(stepParts) == 1 {
			// ...-...
			if strings.Contains(part, "-") {
				rangeParts := strings.Split(part, "-")
				if len(rangeParts) != 2 {
					return fmt.Errorf("invalid %s field", cf.Name)
				}

				start, err := Trans(rangeParts[0])
				if err != nil || start < cf.MinValue || start > cf.MaxValue {
					return fmt.Errorf("invalid %s field", cf.Name)
				}

				end, err := Trans(rangeParts[1])
				if err != nil || end < cf.MinValue || end > cf.MaxValue || end < start {
					return fmt.Errorf("invalid %s field", cf.Name)
				}

				for i := start; i <= end; i++ {
					cf.Value = append(cf.Value, i)
				}
			} else {
				// ...
				value, err := Trans(part)
				if err != nil || value < cf.MinValue || value > cf.MaxValue {
					return fmt.Errorf("invalid %s field", cf.Name)
				}

				cf.Value = append(cf.Value, value)
			}
		} else {
			return fmt.Errorf("invalid %s field", cf.Name)
		}
	}
	return nil
}
