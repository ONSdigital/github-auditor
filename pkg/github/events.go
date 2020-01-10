package github

// MessageForEvent returns a description string with format specifiers for the passed GitHub action.
func MessageForEvent(action string) string {
	eventMap := map[string]string{

		// OAuth events.
		"oauth_application.create": "New OAuth app *%s* was created within organisation *%s* by user *%s*.",

		// Organisation events.
		"org.add_member":                     "User *%s* accepted an invitation to have collaboration access to repo *%s*.",
		"org.block_user":                     "User *%s* was blocked by user *%s*.",
		"org.create":                         "Organisation *%s* was created by user *%s*.",
		"org.disable_saml":                   "SAML was disabled for organisation *%s* by user *%s*.",
		"org.disable_two_factor_requirement": "Two-factor authentication was disabled for organisation *%s* by user *%s*.",
		"org.enable_oauth_app_restrictions":  "OAuth app restrictions were enabled for organisation *%s* by user *%s*.",
		"org.enable_saml":                    "SAML was enabled for organisation *%s* by user *%s*.",
		"org.enable_two_factor_requirement":  "Two-factor authentication was enabled for organisation *%s* by user *%s*.",
		"org.invite_member":                  "New user *%s* was invited to join organisation *%s* by user *%s*.",
		"org.oauth_app_access_approved":      "OAuth app *%s* within organisation *%s* had access approved by user *%s*.",
		"org.oauth_app_access_denied":        "OAuth app *%s* within organisation *%s* had access denied by user *%s*.",
		"org.oauth_app_access_requested":     "Access to OAuth app *%s* within organisation *%s* was requested by user *%s*.",
		"org.remove_billing_manager":         "Billing manager *%s* was removed from organisation *%s* because *%s*.",
		"org.remove_member":                  "User *%s* removed user *%s* from organisation *%s*.",

		// Repo events.
		"repo.access":        "User *%s* changed visibility of repo *%s* to *%s*.",
		"repo.add_member":    "User *%s* accepted invitation to collaborate on repo *%s*.",
		"repo.archived":      "User *%s* archived repo *%s*.",
		"repo.create":        "User *%s* created repo *%s* with visibility *%s*.",
		"repo.destroy":       "User *%s* deleted repo *%s*.",
		"repo.remove_member": "User *%s* was removed as a colloborator from repo *%s*.",

		// Team events.
		"team.add_repository":    "User *%s* gave team *%s* control of repository *%s*.",
		"team.remove_member":     "User *%s* removed user *%s* from team *%s*.",
		"team.remove_repository": "User *%s* removed control from team *%s* of repository *%s*.",
	}

	return eventMap[action]
}
