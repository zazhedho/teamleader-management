package interfacemetric

import domainmetric "teamleader-management/internal/domain/metric"

type RepoMetricInterface interface {
	SaveQuizResults(entries []domainmetric.QuizResult) error
	SaveAppleLogins(entries []domainmetric.AppleLogin) error
	SaveSalesFLP(entries []domainmetric.SalesFLP) error
	SaveApplePoints(entries []domainmetric.ApplePoint) error
	SaveMyHeroPoints(entries []domainmetric.MyHeroPoint) error
	SaveProspects(entries []domainmetric.Prospect) error
}
