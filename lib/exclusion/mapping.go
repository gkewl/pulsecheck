package exclusion

// NamedError holds the error code, short description, and HTTP status
// type NamedExclusion struct {

// }

var exclusions = map[string]string{}

// InitializeExclusions -
func InitializeExclusions() {
	exclusions = make(map[string]string)

	// *** start OIG  ** //
	//Mandatory Exclusions
	exclusions["1128a1"] = "Conviction of program-related crimes. Minimum Period: 5 years"
	exclusions["1128a2"] = "Conviction relating to patient abuse or neglect. Minimum Period: 5 years"
	exclusions["1128a3"] = "Felony conviction relating to health care fraud. Minimum Period: 5 years"
	exclusions["1128a4"] = "Felony conviction relating to controlled substance. Minimum Period: 5 years"
	exclusions["1128c3Gi"] = "Conviction of two mandatory exclusion offenses. Minimum Period: 10 years"
	exclusions["1128c3Gii"] = "Conviction on 3 or more occasions of mandatory exclusion offenses. Permanent Exclusion"
	//Permissive Exclusions
	exclusions["1128b1A"] = "Misdemeanor conviction relating to health care fraud. Baseline Period: 3 years"
	exclusions["1128b1B"] = "Conviction relating to fraud in non- health care programs. Baseline Period: 3"
	exclusions["1128b2"] = "Conviction relating to obstruction of an investigation. Baseline Period: 3 years"
	exclusions["1128b3"] = "Misdemeanor conviction relating to controlled substance. Baseline Period: 3 years"
	exclusions["1128b4"] = "License revocation or suspension. Minimum Period: No less than the period imposed by the state licensing authority."
	exclusions["1128b5"] = "Exclusion or suspension under federal or state health care program. Minimum Period: No less than the period imposed by federal or state health care program."
	exclusions["1128b6"] = "Claims for excessive charges, unnecessary services or services which fail to meet professionally recognized standards of health care, or failure of an HMO to furnish medically necessary services. Minimum Period: 1 year"
	exclusions["1128b7"] = "Fraud, kickbacks, and other prohibited activities. Minimum Period: None"
	exclusions["1128b8"] = "Entities controlled by a sanctioned individual. Minimum Period: Same as length of individual's exclusion."
	exclusions["1128b8A"] = "Entities controlled by a family or household member of an excluded individual and where there has been a transfer of ownership/ control. Minimum Period: Same as length of individual's exclusion."
	exclusions["1128b9"] = "Failure to disclose required information, supply requested information on subcontractors and suppliers; or supply payment information. Minimum Period: None"
	exclusions["1128b10"] = "Failure to disclose required information, supply requested information on subcontractors and suppliers; or supply payment information. Minimum Period: None"
	exclusions["1128b11"] = "Failure to disclose required information, supply requested information on subcontractors and suppliers; or supply payment information. Minimum Period: None"

	exclusions["1128b12"] = "Failure to grant immediate access. Minimum Period: None"
	exclusions["1128b13"] = "Failure to take corrective action. Minimum Period: None"
	exclusions["1128b14"] = "Default on health education loan or scholarship obligations. Minimum Period: Until default has been cured or obligations have been resolved to Public Health Service's (PHS) satisfaction."
	exclusions["1128b15"] = "Individuals controlling a sanctioned entity. Minimum Period: Same period as entity."
	exclusions["1128b16"] = "Making false statement or misrepresentations of material fact. Minimum period: None. The effective date for this new provision is the date of enactment, March 23, 2010."
	exclusions["1156"] = "Failure to meet statutory obligations of practitioners and providers to provide' medically necessary services meeting professionally recognized standards of health care (Quality Improvement Organization (QIO) findings). Minimum Period: 1 year"

	// *** start OIG  ** //

}

// GetDetails -
func GetDetails(exclType string) string {
	return exclusions[exclType]
}
