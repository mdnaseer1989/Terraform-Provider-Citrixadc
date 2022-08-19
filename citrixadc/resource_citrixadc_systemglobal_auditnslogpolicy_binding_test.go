/*
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
        "github.com/citrix/adc-nitro-go/service"
        "github.com/hashicorp/terraform/helper/resource"
        "github.com/hashicorp/terraform/terraform"
        "testing"
)

const testAccSystemglobal_auditnslogpolicy_binding_basic = `

	resource "citrixadc_systemglobal_auditnslogpolicy_binding" "tf_systemglobal_auditnslogpolicy_binding" {
		policyname = "tf_auditnslogpolicy"
		priority   = 50
	}
`

const testAccSystemglobal_auditnslogpolicy_binding_basic_step2 = `
        # Keep the above bound resources without the actual binding to check proper deletion
`

func TestAccSystemglobal_auditnslogpolicy_binding_basic(t *testing.T) {
        resource.Test(t, resource.TestCase{
                PreCheck:     func() { testAccPreCheck(t) },
                Providers:    testAccProviders,
                CheckDestroy: testAccCheckSystemglobal_auditnslogpolicy_bindingDestroy,
                Steps: []resource.TestStep{
                        resource.TestStep{
                                Config: testAccSystemglobal_auditnslogpolicy_binding_basic,
                                Check: resource.ComposeTestCheckFunc(
                                        testAccCheckSystemglobal_auditnslogpolicy_bindingExist("citrixadc_systemglobal_auditnslogpolicy_binding.tf_systemglobal_auditnslogpolicy_binding", nil),
                                ),
                        },
                        resource.TestStep{
                                Config: testAccSystemglobal_auditnslogpolicy_binding_basic_step2,
                                Check: resource.ComposeTestCheckFunc(
                                        testAccCheckSystemglobal_auditnslogpolicy_bindingNotExist("citrixadc_systemglobal_auditnslogpolicy_binding.tf_systemglobal_auditnslogpolicy_binding", "tf_auditnslogpolicy"),
                                ),
                        },
                },
        })
}

func testAccCheckSystemglobal_auditnslogpolicy_bindingExist(n string, id *string) resource.TestCheckFunc {
        return func(s *terraform.State) error {
                rs, ok := s.RootModule().Resources[n]
                if !ok {
                        return fmt.Errorf("Not found: %s", n)
                }

                if rs.Primary.ID == "" {
                        return fmt.Errorf("No systemglobal_auditnslogpolicy_binding id is set")
                }

                if id != nil {
                        if *id != "" && *id != rs.Primary.ID {
                                return fmt.Errorf("Resource ID has changed!")
                        }

                        *id = rs.Primary.ID
                }

                client := testAccProvider.Meta().(*NetScalerNitroClient).client

                policyname := rs.Primary.ID
               
				findParams := service.FindParams{
                        ResourceType:             "systemglobal_auditnslogpolicy_binding",
                        ResourceMissingErrorCode: 258,
                }
                dataArr, err := client.FindResourceArrayWithParams(findParams)

                // Unexpected error
                if err != nil {
                        return err
                }

                // Iterate through results to find the one with the matching policyname
                found := false
                for _, v := range dataArr {
                        if v["policyname"].(string) == policyname {
                                found = true
                                break
                        }
                }

                if !found {
                        return fmt.Errorf("systemglobal_auditnslogpolicy_binding %s not found", n)
                }

                return nil
        }
}

func testAccCheckSystemglobal_auditnslogpolicy_bindingNotExist(n string, id string) resource.TestCheckFunc {
        return func(s *terraform.State) error {
                client := testAccProvider.Meta().(*NetScalerNitroClient).client
                
				policyname := id
                findParams := service.FindParams{
                        ResourceType:             "systemglobal_auditnslogpolicy_binding",
                        ResourceMissingErrorCode: 258,
                }
                dataArr, err := client.FindResourceArrayWithParams(findParams)

                // Unexpected error
                if err != nil {
                        return err
                }

                // Iterate through results to hopefully not find the one with the matching policyname
                found := false
                for _, v := range dataArr {
                        if v["policyname"].(string) == policyname {
                                found = true
                                break
                        }
                }

                if found {
                        return fmt.Errorf("systemglobal_auditnslogpolicy_binding %s was found, but it should have been destroyed", n)
                }

                return nil
        }
}

func testAccCheckSystemglobal_auditnslogpolicy_bindingDestroy(s *terraform.State) error {
        nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client

        for _, rs := range s.RootModule().Resources {
                if rs.Type != "citrixadc_systemglobal_auditnslogpolicy_binding" {
                        continue
                }

                if rs.Primary.ID == "" {
                        return fmt.Errorf("No name is set")
                }

                _, err := nsClient.FindResource(service.Systemglobal_auditnslogpolicy_binding.Type(), rs.Primary.ID)
                if err == nil {
                        return fmt.Errorf("systemglobal_auditnslogpolicy_binding %s still exists", rs.Primary.ID)
                }

        }

        return nil
}