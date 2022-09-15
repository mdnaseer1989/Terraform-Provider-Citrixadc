/*
Copyright 2016 Citrix Systems, Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package citrixadc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

const testAccPcpserver_basic = `

	resource "citrixadc_pcpserver" "tf_pcpserver" {
		name      = "my_pcpserver"
		ipaddress = "10.222.74.185"
		port      = 5351
	}
  
`
const testAccPcpserver_update = `

	resource "citrixadc_pcpserver" "tf_pcpserver" {
		name      = "my_pcpserver"
		ipaddress = "10.222.74.185"
		port      = 5352
	}
  
`


func TestAccPcpserver_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPcpserverDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPcpserver_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPcpserverExist("citrixadc_pcpserver.tf_pcpserver", nil),
					resource.TestCheckResourceAttr("citrixadc_pcpserver.tf_pcpserver", "name", "my_pcpserver"),
					resource.TestCheckResourceAttr("citrixadc_pcpserver.tf_pcpserver", "ipaddress", "10.222.74.185"),
					resource.TestCheckResourceAttr("citrixadc_pcpserver.tf_pcpserver", "port", "5351"),
				),
			},
			resource.TestStep{
				Config: testAccPcpserver_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPcpserverExist("citrixadc_pcpserver.tf_pcpserver", nil),
					resource.TestCheckResourceAttr("citrixadc_pcpserver.tf_pcpserver", "name", "my_pcpserver"),
					resource.TestCheckResourceAttr("citrixadc_pcpserver.tf_pcpserver", "ipaddress", "10.222.74.185"),
					resource.TestCheckResourceAttr("citrixadc_pcpserver.tf_pcpserver", "port", "5352"),
				),
			},
		},
	})
}

func testAccCheckPcpserverExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No pcpserver name is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client
		data, err := nsClient.FindResource("pcpserver", rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("pcpserver %s not found", n)
		}

		return nil
	}
}

func testAccCheckPcpserverDestroy(s *terraform.State) error {
	nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "citrixadc_pcpserver" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := nsClient.FindResource("pcpserver", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("pcpserver %s still exists", rs.Primary.ID)
		}

	}

	return nil
}
