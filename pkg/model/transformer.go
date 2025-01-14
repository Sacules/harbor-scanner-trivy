package model

import (
	"github.com/aquasecurity/harbor-scanner-trivy/pkg/model/harbor"
	"github.com/aquasecurity/harbor-scanner-trivy/pkg/model/trivy"
	log "github.com/sirupsen/logrus"
	"time"
)

type Transformer interface {
	Transform(req harbor.ScanRequest, source trivy.ScanResult) harbor.ScanResult
}

func NewTransformer() Transformer {
	return &transformer{}
}

type transformer struct {
}

func (t *transformer) Transform(req harbor.ScanRequest, source trivy.ScanResult) (target harbor.ScanResult) {
	var vulnerabilities []harbor.VulnerabilityItem

	for _, v := range source.Vulnerabilities {
		vulnerabilities = append(vulnerabilities, harbor.VulnerabilityItem{
			ID:          v.VulnerabilityID,
			Pkg:         v.PkgName,
			Version:     v.InstalledVersion,
			FixVersion:  v.FixedVersion,
			Severity:    t.toHarborSeverity(v.Severity),
			Description: v.Description,
			Links:       v.References,
		})
	}

	target = harbor.ScanResult{
		GeneratedAt: time.Now(),
		Scanner: harbor.Scanner{
			Name:    "Trivy",
			Vendor:  "Aqua Security",
			Version: "0.1.6",
		},
		Artifact:        req.Artifact,
		Severity:        t.toHighestSeverity(source),
		Vulnerabilities: vulnerabilities,
	}
	return
}

func (t *transformer) toHarborSeverity(severity string) harbor.Severity {
	switch severity {
	case "CRITICAL":
		return harbor.SevCritical
	case "HIGH":
		return harbor.SevHigh
	case "MEDIUM":
		return harbor.SevMedium
	case "LOW":
		return harbor.SevLow
	case "UNKNOWN":
		return harbor.SevUnknown
	default:
		log.Printf("Unknown trivy severity %s", severity)
		return harbor.SevUnknown
	}
}

func (t *transformer) toHighestSeverity(sr trivy.ScanResult) (highest harbor.Severity) {
	highest = harbor.SevNone

	for _, vln := range sr.Vulnerabilities {
		sev := t.toHarborSeverity(vln.Severity)
		if sev > highest {
			highest = sev
		}
	}

	return
}
