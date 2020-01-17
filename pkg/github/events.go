package github

// MessageForEvent returns a description string with format specifiers for the passed GitHub action.
// To avoid confusion, the user who initiated an action (the actor) should always appear before the
// user affected by the action (the user).
func MessageForEvent(action string) string {
	eventMap := map[string]string{

		// OAuth events.
		"oauth_application.create": "New OAuth app *%s* was created within organisation *%s* by user *%s*.",

		// Organisation events.
		"org.add_billing_manager":            "%s added %s as billing manager for organisation *%s*.",
		"org.add_member":                     "%s accepted invitation to join organisation *%s*.",
		"org.block_user":                     "%s was blocked by %s in organisation *%s*.",
		"org.create":                         "Organisation *%s* was created by %s.",
		"org.disable_saml":                   "SAML was disabled for organisation *%s* by %s.",
		"org.disable_two_factor_requirement": "Two-factor authentication was disabled for organisation *%s* by %s.",
		"org.enable_oauth_app_restrictions":  "OAuth app restrictions were enabled for organisation *%s* by %s.",
		"org.enable_saml":                    "SAML was enabled for organisation *%s* by %s.",
		"org.enable_two_factor_requirement":  "Two-factor authentication was enabled for organisation *%s* by %s.",
		"org.invite_member":                  "%s invited new %s to join organisation *%s*.",
		"org.oauth_app_access_approved":      "OAuth app *%s* within organisation *%s* had access approved by %s.",
		"org.oauth_app_access_denied":        "OAuth app *%s* within organisation *%s* had access denied by %s.",
		"org.oauth_app_access_requested":     "Access to OAuth app *%s* within organisation *%s* was requested by %s.",
		"org.remove_billing_manager":         "%s removed %s as billing manager from organisation *%s*.",
		"org.remove_member":                  "%s removed %s from organisation *%s*.",
		"org.remove_outside_collaborator":    "%s removed %s as an outside collaborator from organisation *%s*.",
		"org.update_member":                  "%s changed the role of %s from *%s* to *%s* in organisation *%s*.",

		// Repo events.
		"repo.access":               "%s changed the visibility of repo *%s* to *%s*.",
		"repo.add_member":           "%s invited %s to collaborate on repo *%s*.",
		"repo.archived":             "%s archived repo *%s*.",
		"repo.change_merge_setting": "%s changed the merge setting of repo *%s* to *%s*.",
		"repo.create":               "%s created repo *%s* with visibility *%s*.",
		"repo.destroy":              "%s deleted repo *%s*.",
		"repo.remove_member":        "%s removed %s as a collaborator from repo *%s*.",

		// Team events.
		"team.add_member":        "%s added %s to team *%s*.",
		"team.add_repository":    "%s gave team *%s* control of repository *%s*.",
		"team.remove_member":     "%s removed %s from team *%s*.",
		"team.remove_repository": "%s removed control from team *%s* of repository *%s*.",
	}

	return eventMap[action]
}
