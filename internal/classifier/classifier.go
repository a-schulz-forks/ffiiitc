package classifier

import (
	"github.com/go-pkgz/lgr"
	"github.com/navossoc/bayesian"
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
	"log"
	"regexp"
	"slices"
	"strings"
)

// classifier implementation
type TrnClassifier struct {
	Classifier *bayesian.Classifier
	logger     *lgr.Logger
}

type TransactionDataSet [][]string

// init classifier with model file
func NewTrnClassifierFromFile(modelFile string, l *lgr.Logger) (*TrnClassifier, error) {
	cls, err := bayesian.NewClassifierFromFile(modelFile)
	if err != nil {
		return nil, err
	}
	return &TrnClassifier{
		Classifier: cls,
		logger:     l,
	}, nil
}

// NewTrnClassifierWithTraining init classifier with training data set
func NewTrnClassifierWithTraining(catList []string, dataSet TransactionDataSet, l *lgr.Logger) (*TrnClassifier, error) {
	trainingMap := convertDatasetToTrainingMap(dataSet)
	//catList := getCategoriesFromTrainingMap(trainingMap)
	//catList := maps.Keys(trainingMap)
	var classList []bayesian.Class
	for _, str := range catList {
		classList = append(classList, bayesian.Class(str))
	}
	cls := bayesian.NewClassifier(classList...)
	for _, cat := range classList {
		cls.Learn(trainingMap[string(cat)], cat)
	}
	return &TrnClassifier{
		Classifier: cls,
		logger:     l,
	}, nil
}

// save classifier to model file
func (tc *TrnClassifier) SaveClassifierToFile(modelFile string) error {
	err := tc.Classifier.WriteToFile(modelFile)
	return err
}

// perform transaction classification
// in: transaction description
// out: likely transaction category
func (tc *TrnClassifier) ClassifyTransaction(t string) string {
	features := ExtractTransactionFeatures(t)
	_, likely, _ := tc.Classifier.LogScores(features)
	return string(tc.Classifier.Classes[likely])
}

// function to get category and list of
// unique features from line of transaction data set
// in: [cat, trn description]
// out: cat, [features...]
func getCategoryAndFeatures(data []string) (string, []string) {
	category := data[0]
	var features = ExtractTransactionFeatures(data[1])
	return category, features

}

// get slice of categories from training map
func getCategoriesFromTrainingMap(training map[string][]string) []bayesian.Class {
	var result []bayesian.Class
	for key := range training {
		result = append(result, bayesian.Class(key))
	}
	return result
}

// checks if string is pure number: int or float
func isStringNumeric(s string) bool {
	numericPattern := `^-?\d+(\.\d+)?$`
	match, err := regexp.MatchString(numericPattern, s)
	return err == nil && match
}

// build training map from transactions data set
// in: [ [cat, trn description], [cat, trn description]... ]
// out: map[Category] = [feature1, feature2, ...]
func convertDatasetToTrainingMap(dataSet TransactionDataSet) map[string][]string {
	resultMap := make(map[string][]string)
	var features []string
	var category string
	for _, line := range dataSet {
		category, features = getCategoryAndFeatures(line)
		_, exist := resultMap[category]
		if exist {
			resultMap[category] = append(resultMap[category], features...)
		} else {
			resultMap[category] = features
		}
	}
	return resultMap
}

// ExtractTransactionFeatures extract unique words from transaction description that are not numeric
func ExtractTransactionFeatures(transaction string) []string {
	var transFeatures []string
	configFile, err := tokenizer.CachedPath("bert-base-german-cased", "tokenizer.json")
	if err != nil {
		panic(err)
	}

	tk, err := pretrained.FromFile(configFile)
	if err != nil {
		panic(err)
	}
	en, err := tk.EncodeSingle(transaction)
	if err != nil {
		log.Fatal(err)
	}
	features := en.Tokens
	for _, feature := range features {
		feature = strings.TrimPrefix(feature, "##")
		// Find the index of "DATUM" in the string
		index := strings.Index(feature, "DATUM")

		// If "DATUM" is found, trim everything from that point onward
		if index != -1 {
			feature = strings.TrimSpace(feature[:index])
		}

		if (len(feature) > 1) && (!slices.Contains(transFeatures, feature)) && (!isStringNumeric(feature)) {
			transFeatures = append(transFeatures, feature)
		}
	}
	return transFeatures
}
