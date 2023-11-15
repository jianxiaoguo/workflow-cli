package cmd

import (
	"fmt"
	"os"

	"github.com/drycc/controller-sdk-go/routes"
	"sigs.k8s.io/yaml"
)

// RoutesCreate create a route to an app.
func (d *DryccCmd) RoutesCreate(appID, name string, procType string, kind string, port int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	d.Printf("Adding route %s to %s... ", name, appID)

	quit := progress(d.WOut)
	err = routes.New(s.Client, appID, name, procType, kind, port)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// RoutesList lists routes for the app
func (d *DryccCmd) RoutesList(appID string, results int) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	if results == defaultLimit {
		results = s.Limit
	}

	routes, count, err := routes.List(s.Client, appID, results)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	if count == 0 {
		d.Println(fmt.Sprintf("No routes found in %s app.", appID))
	} else {
		table := d.getDefaultFormatTable([]string{"NAME", "OWNER", "TYPE", "KIND", "SERVICE-PORT", "GATEWAY", "LISTENER-PORT"})
		for _, route := range routes {
			if len(route.ParentRefs) > 0 {
				for _, gateway := range route.ParentRefs {
					table.Append([]string{
						route.Name,
						route.Owner,
						route.Type,
						route.Kind,
						fmt.Sprint(route.Port),
						gateway.Name,
						fmt.Sprint(gateway.Port),
					})
				}
			} else {
				table.Append([]string{
					route.Name,
					route.Owner,
					route.Type,
					route.Kind,
					fmt.Sprint(route.Port),
					"",
					"",
				})
			}
		}
		table.Render()
	}
	return nil
}

// RoutesAttach bind a route to gateway.
func (d *DryccCmd) RoutesAttach(appID, name string, port int, gateway string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	d.Printf("Attaching route %s to gateway %s... ", name, gateway)

	quit := progress(d.WOut)
	err = routes.AttachGateway(s.Client, appID, name, port, gateway)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// RoutesDetach bind a route to gateway.
func (d *DryccCmd) RoutesDetach(appID, name string, port int, gateway string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	d.Printf("Detaching route %s to gateway %s... ", name, gateway)

	quit := progress(d.WOut)
	err = routes.DetachGateway(s.Client, appID, name, port, gateway)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}

// RoutesGet get rule of route for the app
func (d *DryccCmd) RoutesGet(appID string, name string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}

	route, err := routes.GetRule(s.Client, appID, name)
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	var rules []byte
	rules, err = yaml.JSONToYAML([]byte(route))
	if err != nil {
		return err
	}
	d.Println(string(rules))
	return nil
}

// RoutesSet set rule of route for the app
func (d *DryccCmd) RoutesSet(appID string, name string, ruleFile string) error {
	s, appID, err := load(d.ConfigFile, appID)
	if err != nil {
		return err
	}

	var contents []byte
	if _, err := os.Stat(ruleFile); err != nil {
		return err
	}
	contents, err = os.ReadFile(ruleFile)
	if err != nil {
		return err
	}
	rules, err := yaml.YAMLToJSON(contents)
	if err != nil {
		return err
	}
	d.Print("Applying rules... ")
	quit := progress(d.WOut)
	err = routes.SetRule(s.Client, appID, name, string(rules))
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}
	d.Println("done")
	return nil
}

// RoutesRemove removes a route registered with an app.
func (d *DryccCmd) RoutesRemove(appID, name string) error {
	s, appID, err := load(d.ConfigFile, appID)

	if err != nil {
		return err
	}
	d.Printf("Removing route %s to %s... ", name, appID)

	quit := progress(d.WOut)
	err = routes.Delete(s.Client, appID, name)
	quit <- true
	<-quit
	if d.checkAPICompatibility(s.Client, err) != nil {
		return err
	}

	d.Println("done")
	return nil
}
