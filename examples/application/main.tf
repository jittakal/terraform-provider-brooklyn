provider "brooklyn" {  
  access_key = "<access-key>"
  secret_key = "<secret-key>"
  endpoint_url = "<endpoint-url>"
}

resource "brooklyn_application" "exampleapplication" {
    application_name = "example-application" 
    location_name = "AWS Ireland"
    catalog_id = ""
    catalog_version = ""
}