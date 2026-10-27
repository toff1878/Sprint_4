package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	splitString := strings.Split(data, ",")
	if len(splitString) != 2 {
		return 0, 0, errors.New("некорректное количество входных параметров")
	}
	stepsNumber, err := strconv.Atoi(splitString[0])
	if err != nil {
		return 0, 0, err
	}
	if stepsNumber <= 0 {
		return 0, 0, errors.New("число шагов не может быть отрицательным или равняться нулю")
	}
	walkDuration, err := time.ParseDuration(splitString[1])
	if err != nil {
		return 0, 0, err
	}

	if walkDuration <= 0 {
		return 0, 0, errors.New("продолжительность прогулки не может быть отрицательной или равняться 0")
	}

	return stepsNumber, walkDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	stepsNumber, walkDuration, err := parsePackage(data)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	if stepsNumber <= 0 {
		return ""
	}

	distanceKm := (float64(stepsNumber) * stepLength) / mInKm

	calories, err := spentcalories.WalkingSpentCalories(stepsNumber, weight, height, walkDuration)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		stepsNumber,
		distanceKm,
		calories,
	)

}
