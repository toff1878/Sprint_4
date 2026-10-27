package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	splitString := strings.Split(data, ",")
	if len(splitString) != 3 {
		return 0, "", 0, errors.New("некорректное количество входных параметров")
	}
	stepsNumber, err := strconv.Atoi(splitString[0])
	if err != nil {
		return 0, "", 0, err
	}
	if stepsNumber <= 0 {
		return 0, "", 0, errors.New("число шагов не может быть отрицательным или равняться нулю")
	}
	trainingType := splitString[1]

	walkDuration, err := time.ParseDuration(splitString[2])
	if err != nil {
		return 0, "", 0, err
	}

	if walkDuration <= 0 {
		return 0, "", 0, errors.New("продолжительность тренировки не может быть отрицательной или равняться 0")
	}
	return stepsNumber, trainingType, walkDuration, nil
}

func distance(steps int, height float64) float64 {

	stepLength := stepLengthCoefficient * height

	meters := float64(steps) * stepLength

	return meters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distanceKm := distance(steps, height)

	return distanceKm / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	stepsNumber, trainingType, walkDuration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	walkCaloriesSpent, err := WalkingSpentCalories(stepsNumber, weight, height, walkDuration)
	if err != nil {
		log.Println(err)
		return "", err
	}
	runCaloriesSpent, err := RunningSpentCalories(stepsNumber, weight, height, walkDuration)
	if err != nil {
		log.Println(err)
		return "", err
	}

	switch trainingType {
	case "Ходьба":

		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			trainingType,
			walkDuration.Hours(),
			distance(stepsNumber, height),
			meanSpeed(stepsNumber, height, walkDuration),
			walkCaloriesSpent,
		), nil
	case "Бег":
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			trainingType,
			walkDuration.Hours(),
			distance(stepsNumber, height),
			meanSpeed(stepsNumber, height, walkDuration),
			runCaloriesSpent,
		), nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("число шагов не может быть отрицательным или равняться нулю")
	}

	if weight <= 15 {
		return 0, errors.New("вес пользователя не может быть меньше 15 кг")
	}

	if height <= 0.5 {
		return 0, errors.New("рост пользователя не может быть меньше 50 см")
	}

	if duration <= 0 {
		return 0, errors.New("продолжительность тренировки не может быть отрицательной или равняться 0")
	}

	speed := meanSpeed(steps, height, duration)

	return (weight * speed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("число шагов не может быть отрицательным или равняться нулю")
	}

	if weight <= 15 {
		return 0, errors.New("вес пользователя не может быть меньше 15 кг")
	}

	if height <= 0.5 {
		return 0, errors.New("рост пользователя не может быть меньше 50 см")
	}

	if duration <= 0 {
		return 0, errors.New("продолжительность тренировки не может быть отрицательной или равняться 0")
	}

	speed := meanSpeed(steps, height, duration)

	return ((weight * speed * duration.Minutes()) / minInH) * walkingCaloriesCoefficient, nil
}
