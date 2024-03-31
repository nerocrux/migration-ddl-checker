package ddl

type DDL string

const (
	CreateDDL DDL = "CREATE"
	DropDDL   DDL = "DROP"
	AlterDDL  DDL = "ALTER"
)

func FromConfig(hazardousDDLs []string) []DDL {
	var ddls []DDL
	for _, input := range hazardousDDLs {
		ddls = append(ddls, DDL(input))
	}
	return ddls
}
