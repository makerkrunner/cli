package v7action

// Organization represents a V7 actor organization.
type OrganizationQuota struct {
	// GUID is the unique organization identifier.
	GUID string
	// Name is the name of the organization.
	Name string
}

// CreateOrganization creates a new organization with the given name
func (actor Actor) CreateOrganizationQuota(orgQuotaName string) (OrganizationQuota, Warnings, error) {
	allWarnings := Warnings{}

	// organization, apiWarnings, err := actor.CloudControllerClient.CreateOrganization(orgName)
	// allWarnings = append(allWarnings, apiWarnings...)

	return OrganizationQuota{}, allWarnings, nil
}
