package report

import "fmt"

func PrintSummary(staleFlags []string, removedFlags []string) {
	if len(staleFlags) == 0 {
		fmt.Println("No stale flags found")
		return
	}
	fmt.Printf("Unleash Checker found %d stale flags\n", len(staleFlags))
	if len(removedFlags) > 0 {
		fmt.Println("The following flags were removed from files due to being stale:")
		for _, flag := range removedFlags {
			fmt.Printf(" - %s\n", flag)
		}
		fmt.Println("")
	}
	unfoundFlags := difference(staleFlags, removedFlags)
	if len(unfoundFlags) > 0 {	
		fmt.Println("The following flags were not found in the specified directory:")
		for _, flag := range unfoundFlags {
			fmt.Printf(" - %s\n", flag)
		}
		fmt.Println("")
	}
	fmt.Println("Please review the changes and commit them to your repository.\nIf you still want to use these flags, consider changing the flag type:\nhttps://docs.getunleash.io/reference/technical-debt")
}

func difference(sliceA, sliceB []string) []string {
	var diff []string

	for _, a := range sliceA {
		found := false
		for _, b := range sliceB {
			if a == b {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, a)
		}
	}

	return diff
}
