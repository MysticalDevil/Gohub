package config

import "gohub/pkg/config"

func init() {
	config.Add("paging", func() map[string]any {
		return map[string]any{
			// Default number of entries per page
			"perpage": 10,

			// The parameter in the URL to distinguish how many pages
			// If this value is changed, the request validation rule must be changed as well
			"url_query_page": "page",

			// The parameters in the URL to distinguish sorting (using id or other)
			// If this value is changed, the request validation rule must be changed as well
			"url_query_sort": "sort",

			// The parameters in the URL to distinguish sorting rules (forward or reverse order)
			// If this value is changed, the request validation rule must be changed as well
			"url_query_order": "order",

			// The parameter in the URL to distinguish number of entries per page
			// If this value is changed, the request validation rule must be changed as well
			"url_query_per_page": "per_page",
		}
	})
}
