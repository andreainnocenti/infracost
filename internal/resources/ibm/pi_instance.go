package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/resources"
	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// PiInstance struct represents <TODO: cloud service short description>.
//
// <TODO: Add any important information about the resource and links to the
// pricing pages or documentation that might be useful to developers in the future, e.g:>
//
// Resource information: https://cloud.ibm.com/<PATH/TO/RESOURCE>/
// Pricing information: https://cloud.ibm.com/<PATH/TO/PRICING>/
type PiInstance struct {
	Address       string
	Region        string
	ProcessorMode string
	SystemType    string
	StorageType   string
	Memory        float64
	Cpus          float64

	Storage              *float64 `infracost_usage:"storage"`
	MonthlyInstanceHours *float64 `infracost_usage:"monthly_instance_hours"`
}

// PiInstanceUsageSchema defines a list which represents the usage schema of PiInstance.
var PiInstanceUsageSchema = []*schema.UsageItem{
	{Key: "storage", DefaultValue: 0, ValueType: schema.Float64},
	{Key: "monthly_instance_hours", DefaultValue: 0, ValueType: schema.Float64},
}

// PopulateUsage parses the u schema.UsageData into the PiInstance.
// It uses the `infracost_usage` struct tags to populate data into the PiInstance.
func (r *PiInstance) PopulateUsage(u *schema.UsageData) {
	resources.PopulateArgsWithUsage(r, u)
}

// BuildResource builds a schema.Resource from a valid PiInstance struct.
// This method is called after the resource is initialised by an IaC provider.
// See providers folder for more information.
func (r *PiInstance) BuildResource() *schema.Resource {
	costComponents := []*schema.CostComponent{
		r.piInstanceCoresCostComponent(),
		r.piInstanceMemoryCostComponent(),
		r.piInstanceStorageCostComponent(),
	}

	return &schema.Resource{
		Name:           r.Address,
		UsageSchema:    PiInstanceUsageSchema,
		CostComponents: costComponents,
	}
}

func (r *PiInstance) piInstanceCoresCostComponent() *schema.CostComponent {

	var q *decimal.Decimal

	if r.MonthlyInstanceHours != nil {
		q = decimalPtr(decimal.NewFromFloat(r.Cpus * (*r.MonthlyInstanceHours)))
	}

	const s922 string = "s922"
	const e980 string = "e980"
	const e1080 string = "e1080"

	unit := ""

	if r.ProcessorMode == "shared" {
		if r.SystemType == s922 {
			unit = "SOS_VIRTUAL_PROCESSOR_CORE_HOURS"
		} else if r.SystemType == e980 {
			unit = "ESS_VIRTUAL_PROCESSOR_CORE_HOURS"
		} else if r.SystemType == e1080 {
			unit = "PTEN_ESS_VIRTUAL_PROCESSOR_CORE_HRS"
		}
	} else if r.ProcessorMode == "dedicated" {
		if r.SystemType == s922 {
			unit = "SOD_VIRTUAL_PROCESSOR_CORE_HOURS"
		} else if r.SystemType == e980 {
			unit = "EDD_VIRTUAL_PROCESSOR_CORE_HOURS"
		} else if r.SystemType == e1080 {
			unit = "PTEN_EDD_VIRTUAL_PROCESSOR_CORE_HRS"
		}
	} else if r.ProcessorMode == "capped" {
		if r.SystemType == s922 {
			unit = "SOC_VIRTUAL_PROCESSOR_CORE_HOURS"
		} else if r.SystemType == e980 {
			unit = "ECC_VIRTUAL_PROCESSOR_CORE_HOURS"
		} else if r.SystemType == e1080 {
			unit = "PTEN_ECC_VIRTUAL_PROCESSOR_CORE_HRS"
		}
	}

	return &schema.CostComponent{
		Name:            "Cores",
		Unit:            "Core Hours",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceMemoryCostComponent() *schema.CostComponent {

	var q *decimal.Decimal

	if r.MonthlyInstanceHours != nil {
		q = decimalPtr(decimal.NewFromFloat(r.Memory * (*r.MonthlyInstanceHours)))
	}

	unit := "MS_GIGABYTE_HOURS"

	return &schema.CostComponent{
		Name:            "Memory",
		Unit:            "GB Hours",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}

func (r *PiInstance) piInstanceStorageCostComponent() *schema.CostComponent {

	var q *decimal.Decimal

	if r.Storage != nil {
		q = decimalPtr(decimal.NewFromFloat((*r.Storage) * (*r.MonthlyInstanceHours)))
	}

	unit := ""

	if r.StorageType == "tier1" {
		unit = "TIER_ONE_STORAGE_GIGABYTE_HOURS"
	} else if r.StorageType == "tier3" {
		unit = "TIER_THREE_STORAGE_GIGABYTE_HOURS"
	}

	return &schema.CostComponent{
		Name:            fmt.Sprintf("Storage - %s", r.StorageType),
		Unit:            "GB Hours",
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("ibm"),
			Region:        strPtr(r.Region),
			ProductFamily: strPtr("service"),
			Service:       strPtr("power-iaas"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: strPtr("power-virtual-server-group")},
				{Key: "planType", Value: strPtr("Paid")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(unit),
		},
	}
}
