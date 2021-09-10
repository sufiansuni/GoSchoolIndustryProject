package health

// Calculate BMI from height(cm) and weight(kg) inputs.
// Body mass index (BMI) is a measure of body fat based on height and weight that applies to adult men and women.
// https://www.nhlbi.nih.gov/health/educational/lose_wt/BMI/bmicalc.htm
func CalculateBMI(height float64, weight float64) float64 {
	// Formula: weight (kg) / [height (m)]2
	// The formula for BMI is weight in kilograms divided by height in meters squared.
	// If height has been measured in centimeters, divide by 100 to convert this to meters.
	height = height / 100
	BMI := weight / (height * height)
	return BMI
}

// Determines BMI catergory from BMI value.
// https://theconversation.com/body-mass-index-may-not-be-the-best-indicator-of-our-health-how-can-we-improve-it-143155
func BMICategory(BMI float64) string {
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
