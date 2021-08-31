package main

import "fmt"

// Calculate BMI from height(cm) and weight(kg) inputs.
// Body mass index (BMI) is a measure of body fat based on height and weight that applies to adult men and women.
// https://www.nhlbi.nih.gov/health/educational/lose_wt/BMI/bmicalc.htm
func calculateBMI(height float64, weight float64) float64 {
	// Formula: weight (kg) / [height (m)]2
	// The formula for BMI is weight in kilograms divided by height in meters squared.
	// If height has been measured in centimeters, divide by 100 to convert this to meters.
	height = height / 100
	BMI := weight / (height * height)
	return BMI
}

// Determines BMI catergory from BMI value.
// https://theconversation.com/body-mass-index-may-not-be-the-best-indicator-of-our-health-how-can-we-improve-it-143155
func bmiCategory(BMI float64) string {
	if BMI < 18.5 {
		return "Underweight"
	} else if BMI >= 18.5 && BMI < 25.0 {
		return "Healthy"
	} else if BMI >= 25.0 && BMI < 30.0 {
		return "Overweight"
	} else if BMI >= 30.0 && BMI < 35 {
		return "Obese (Class 1)"
	} else if BMI >= 35.0 && BMI < 40 {
		return "Obese (Class 2)"
	} else {
		return "Obese (Class 3)"
	}
}

// Calculate BMR from gender, height(cm), weight(kg) and age(years) inputs.
// BMR result is in calories per day.
// https://www.calculator.net/bmr-calculator.html
func calculateBMR(gender string, height float64, weight float64, age float64) float64 {
	// Mifflin-St Jeor Equation:
	// W is body weight in kg
	// H is body height in cm
	// A is age
	switch gender {
	case "Male":
		// For men:
		// BMR = 10W + 6.25H - 5A + 5
		return (10 * weight) + (6.25 * height) - (5 * age) + 5
	case "Female":
		// For women:
		// BMR = 10W + 6.25H - 5A - 161
		return (10 * weight) + (6.25 * height) - (5 * age) - 161
	default:
		fmt.Println("Invalid Gender Input when calculating BMR")
		return 0
	}
}

// https://www.omnicalculator.com/health/bmr-harris-benedict-equation
// To determine your total daily calorie needs, multiply your BMR by the appropriate activity factor, as follows:

// Sedentary (little or no exercise) : Calorie-Calculation = BMR x 1.2
// Lightly active (light exercise/sports 1-3 days/week) : Calorie-Calculation = BMR x 1.375
// Moderately active (moderate exercise/sports 3-5 days/week) : Calorie-Calculation = BMR x 1.55
// Very active (hard exercise/sports 6-7 days a week) : Calorie-Calculation = BMR x 1.725
// If you are extra active (very hard exercise/sports & a physical job) : Calorie-Calculation = BMR x 1.9

// Calculate recommended daily calories based on BMR and Activity level inputs.
func calculateDailyCalories(BMR float64, activity int) float64 {
	switch activity {
	case 1:
		return BMR * 1.2
	case 2:
		return BMR * 1.375
	case 3:
		return BMR * 1.55
	case 4:
		return BMR * 1.725
	case 5:
		return BMR * 1.9
	default:
		fmt.Println("Activity Level not 1-5 when calculating recommended daily calories")
		return BMR
	}
}
