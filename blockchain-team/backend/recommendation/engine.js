/**
 * Simulates the Machine Learning Model output.
 * Uses a STATIC timestamp so the generated hash remains identical 
 * across different API calls, ensuring a smooth presentation.
 */
function getMockRecommendation(patientData) {
    const { bmi, diabetes, heartDisease, kidneyDisease } = patientData;
    
    let suitabilityScore = 75; 
    let recommendedDrugs = [];
    let sideEffects = ["Nausea", "Mild Headache"];

    // Simulated ML Logic based on disease profiles
    if (diabetes && !heartDisease && !kidneyDisease) {
        recommendedDrugs = ["Wegovy", "Qsymia"];
        suitabilityScore = 88;
        sideEffects.push("Hypoglycemia risk when paired with insulin");
    } else {
        recommendedDrugs = ["Wegovy", "Zepbound", "Saxenda"];
        suitabilityScore = 90;
    }

    if (bmi && bmi > 40) { 
        suitabilityScore -= 5; 
    }

    return {
        diseaseContext: {
            obesity: true,
            diabetes: !!diabetes,
            heartDisease: !!heartDisease,
            kidneyDisease: !!kidneyDisease
        },
        recommendedDrugs,
        suitabilityScore,
        sideEffects,
        // STATIC TIMESTAMP FOR DEMO STABILITY
        timestamp: "2026-06-06T12:00:00.000Z" 
    };
}

module.exports = { getMockRecommendation };