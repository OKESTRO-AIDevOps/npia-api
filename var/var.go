package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	. "github.com/OKESTRO-AIDevOps/npia-api/pkg/apistandard"

	pkgutils "github.com/OKESTRO-AIDevOps/npia-api/pkg/utils"

	goya "github.com/goccy/go-yaml"

	"strconv"

	bsrc "github.com/OKESTRO-AIDevOps/npia-api/pkg/builtinresource"
	//"github.com/fatih/color"
)

type AppOrigin struct {
	KCFG_PATH string
	MAIN_NS   string
	RECORDS   []RecordInfo
	REPOS     []RepoInfo
	REGS      []RegInfo
}

type RecordInfo struct {
	NS        string
	REPO_ADDR string
	REG_ADDR  string
}

type RepoInfo struct {
	REPO_ADDR string
	REPO_ID   string
	REPO_PW   string
}

type RegInfo struct {
	REG_ADDR string
	REG_ID   string
	REG_PW   string
}

func callApiDef() {

	ASgi.PrintPrettyDefinition()

}

func callApiDefStructure() {

	ASgi.PrintRawDefinition()

}

func sliceTest() {

	ret := pkgutils.InsertToSliceByIndex[string]([]string{"b", "c", "d"}, 0, "a")

	fmt.Println(ret)
}

func writeToAdmOrigin() {

	var test_ao AppOrigin

	var test_ri RecordInfo

	var test_rep RepoInfo

	var test_reg RegInfo

	test_ao.RECORDS = append(test_ao.RECORDS, test_ri)

	test_ao.REPOS = append(test_ao.REPOS, test_rep)

	test_ao.REGS = append(test_ao.REGS, test_reg)

	file_byte, _ := json.Marshal(test_ao)

	_ = os.WriteFile("testadmorigin.json", file_byte, 0644)

}

func yamlLoad(file_path string) {

	file_byte, _ := os.ReadFile(file_path)

	file_list := strings.Split(string(file_byte), "---")

	for _, yaml_file := range file_list {

		readFromYAML(yaml_file, "$.spec.ports[0].port")

	}

}

func readFromYAML(yaml_file string, yaml_path string) {

	ypath, _ := goya.PathString(yaml_path)

	var value int

	_ = ypath.Read(strings.NewReader(yaml_file), &value)

	fmt.Println(value)

}

func komposeTest() {

	cmd := exec.Command("../lib/bin/kompose", "convert", "-f", "../lib/bin/docker-compose.yaml", "--stdout")

	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var test_arr []interface{}

	yaml_str := string(out)

	yaml_path_items := "$.items"

	ypath, _ := goya.PathString(yaml_path_items)

	err = ypath.Read(strings.NewReader(yaml_str), &test_arr)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var to_file_list [][]byte

	var to_file []byte

	for _, val := range test_arr {

		yaml_if := make(map[interface{}]interface{})

		resource_b, err := goya.Marshal(val)

		err = goya.Unmarshal(resource_b, &yaml_if)

		if err != nil {

			fmt.Println(err.Error())
			return
		}

		if yaml_if["kind"] == "Deployment" {

			image_pull_secrets := make([]map[string]string, 0)

			value := map[string]string{
				"name": "docker-secret",
			}

			image_pull_secrets = append(image_pull_secrets, value)

			yaml_if["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["imgaePullSecrets"] = image_pull_secrets

			c_count := len(yaml_if["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{}))

			for j := 0; j < c_count; j++ {

				prefix := "damn/go_"

				prefix += yaml_if["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[j].(map[string]interface{})["image"].(string)

				yaml_if["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[j].(map[string]interface{})["image"] = prefix

				yaml_if["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[j].(map[string]interface{})["imagePullPolicy"] = "Always"
			}
		}

		result_b, err := goya.Marshal(yaml_if)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		to_file_list = append(to_file_list, result_b)

	}

	for i := 0; i < len(to_file_list); i++ {

		to_file = append(to_file, []byte("---\n")...)

		to_file = append(to_file, to_file_list[i]...)

	}

	err = os.WriteFile("done_question_mark.yaml", to_file, 0644)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func dockercomposeyamlTest() {

	file_byte, err := os.ReadFile("../lib/bin/docker-compose.yaml")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	yaml_if := make(map[interface{}]interface{})

	err = goya.Unmarshal(file_byte, &yaml_if)

	if err != nil {

		fmt.Println(err.Error())
		return
	}

	fmt.Println(len(yaml_if["services"].(map[string]interface{})))

	delete(yaml_if["services"].(map[string]interface{})["tgdb"].(map[string]interface{}), "ports")
	delete(yaml_if["services"].(map[string]interface{})["tgdb"].(map[string]interface{}), "volumes")

	tofile, err := goya.Marshal(yaml_if)

	_ = os.WriteFile("done_qm.yaml", tofile, 0644)

}

func delresourceTest() {

	file_byte, err := os.ReadFile("./done_question_mark.yaml")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var kill_doc [][]byte

	var kill_file []byte

	nm_found := 0

	resourcenm := "tgtraffic"

	file_str := string(file_byte)

	file_str_list := strings.Split(file_str, "---\n")

	for _, content := range file_str_list {

		if content == "\n" || content == "" {
			continue
		}

		yaml_if := make(map[interface{}]interface{})

		// c_byte, err := goya.Marshal(content)

		c_byte := []byte(content)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = goya.Unmarshal(c_byte, &yaml_if)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if yaml_if["metadata"].(map[string]interface{})["name"] == resourcenm {

			nm_found = 1

			b_tmp, err := goya.Marshal(yaml_if)

			if err != nil {

				fmt.Println()
				return

			}

			kill_doc = append(kill_doc, b_tmp)
		}

	}

	if nm_found == 0 {
		fmt.Println("matching name not found")
		return
	}

	for _, res_content := range kill_doc {

		kill_file = append(kill_file, []byte("---\n")...)
		kill_file = append(kill_file, res_content...)

	}

	err = os.WriteFile("done_question_mark2.yaml", kill_file, 0644)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func hpaTest() (string, error) {

	resource_key := "Deployment"

	resourcenm := "tgdb"

	src_found := 0

	file_byte, err := os.ReadFile("done_question_mark.yaml")

	if err != nil {
		return "", fmt.Errorf("failed to create hpa src: %s", err.Error())
	}

	cmd := exec.Command("kubectl", "get", "nodes", "-o", "yaml")

	out, err := cmd.Output()

	var yaml_items []interface{}

	yaml_str := string(out)

	yaml_path_items := "$.items"

	ypath, _ := goya.PathString(yaml_path_items)

	err = ypath.Read(strings.NewReader(yaml_str), &yaml_items)

	if err != nil {
		return "", fmt.Errorf("failed to create hpa src: %s", err.Error())
	}

	min_repl := 0

	max_repl := 0

	prev_top := 0

	pods := 0

	for _, val := range yaml_items {

		yaml_if := make(map[interface{}]interface{})

		resource_b, err := goya.Marshal(val)

		err = goya.Unmarshal(resource_b, &yaml_if)

		if err != nil {
			return "", fmt.Errorf("failed to create hpa src: %s", err.Error())
		}

		pods_str := yaml_if["status"].(map[string]interface{})["allocatable"].(map[string]interface{})["pods"].(string)

		pods, err = strconv.Atoi(pods_str)

		if err != nil {
			return "", fmt.Errorf("failed to create hpa src: %s", err.Error())
		}

		if pods > prev_top {
			prev_top = pods
		}

	}

	pods = prev_top

	min_repl = int(float64(pods) * 0.02)

	max_repl = int(float64(pods) * 0.1)

	head_metadataName := "hpa-deployment-" + resourcenm
	apiVersion := ""
	kind := resource_key
	metadata_name := resourcenm

	var out_hpa bsrc.HorizontalPodAutoscaler

	file_str := string(file_byte)

	file_str_list := strings.Split(file_str, "---\n")

	for _, content := range file_str_list {

		if content == "\n" || content == "" {
			continue
		}

		yaml_if := make(map[interface{}]interface{})

		c_byte := []byte(content)

		if err != nil {
			return "", fmt.Errorf("failed to create hpa src: %s", err.Error())
		}

		err = goya.Unmarshal(c_byte, &yaml_if)

		if err != nil {
			return "", fmt.Errorf("failed to create hpa src: %s", err.Error())
		}

		if yaml_if["kind"] == resource_key && yaml_if["metadata"].(map[string]interface{})["name"] == resourcenm {

			src_found = 1
			apiVersion = yaml_if["apiVersion"].(string)
			kind = yaml_if["kind"].(string)
			metadata_name = yaml_if["metadata"].(map[string]interface{})["name"].(string)

			break
		}

	}

	if src_found == 0 {
		return "", fmt.Errorf("failed to create hpa src: %s", "matching key not found")
	}

	if min_repl == 0 {
		min_repl = 1
	}

	if min_repl > max_repl {
		max_repl = min_repl
	}

	out_hpa.APIVersion = "autoscaling/v1"
	out_hpa.Kind = "HorizontalPodAutoscaler"
	out_hpa.Metadata.Name = head_metadataName
	out_hpa.Spec.ScaleTargetRef.APIVersion = apiVersion
	out_hpa.Spec.ScaleTargetRef.Kind = kind
	out_hpa.Spec.ScaleTargetRef.Name = metadata_name
	out_hpa.Spec.MinReplicas = min_repl
	out_hpa.Spec.MaxReplicas = max_repl
	out_hpa.Spec.TargetCPUUtilizationPercentage = 50

	out_hpa_byte, err := goya.Marshal(out_hpa)

	if err != nil {
		return "", fmt.Errorf("failed to create hpa src: %s", err.Error())
	}

	err = os.WriteFile("done_question_mark_hpa.yaml", out_hpa_byte, 0644)

	if err != nil {
		return "", fmt.Errorf("failed to create hpa src: %s", err.Error())
	}

	return ".usr/hpa_src.yaml", nil

}

func qosTest() (string, error) {

	resource_key := "Deployment"

	resourcenm := "tgdb"

	src_found := 0

	file_byte, err := os.ReadFile("done_question_mark.yaml")

	if err != nil {
		return "", fmt.Errorf("failed to create qos src: %s", err.Error())
	}

	cmd := exec.Command("kubectl", "get", "nodes", "-o", "yaml")

	out, err := cmd.Output()

	var yaml_items []interface{}

	yaml_str := string(out)

	yaml_path_items := "$.items"

	ypath, _ := goya.PathString(yaml_path_items)

	err = ypath.Read(strings.NewReader(yaml_str), &yaml_items)

	if err != nil {
		return "", fmt.Errorf("failed to create qos src: %s", err.Error())
	}

	polled_node_index := 0

	prev_top := 0

	pods := 0

	for i, val := range yaml_items {

		yaml_if := make(map[interface{}]interface{})

		resource_b, err := goya.Marshal(val)

		err = goya.Unmarshal(resource_b, &yaml_if)

		if err != nil {
			return "", fmt.Errorf("failed to create qos src: %s", err.Error())
		}

		pods_str := yaml_if["status"].(map[string]interface{})["allocatable"].(map[string]interface{})["pods"].(string)

		pods, err = strconv.Atoi(pods_str)

		if err != nil {
			return "", fmt.Errorf("failed to create qos src: %s", err.Error())
		}

		if pods > prev_top {
			prev_top = pods
			polled_node_index = i

		}

	}

	pods = prev_top

	polled_cpu := yaml_items[polled_node_index].(map[string]interface{})["status"].(map[string]interface{})["allocatable"].(map[string]interface{})["cpu"].(string)

	polled_mem := yaml_items[polled_node_index].(map[string]interface{})["status"].(map[string]interface{})["allocatable"].(map[string]interface{})["memory"].(string)

	pods_float := float64(pods)

	polled_cpu_float, err := strconv.ParseFloat(polled_cpu, 64)

	if err != nil {
		return "", fmt.Errorf("failed to create qos src: %s", err.Error())
	}

	polled_mem_str := strings.ReplaceAll(polled_mem, "Ki", "")

	polled_mem_float, err := strconv.ParseFloat(polled_mem_str, 64)

	if err != nil {
		return "", fmt.Errorf("failed to create qos src: %s", err.Error())
	}

	cpu_limit_per_node := (polled_cpu_float / pods_float) * 8.0

	mem_limit_per_node := (polled_mem_float / pods_float) * 16.0

	qos_cpu_high := strconv.FormatFloat(cpu_limit_per_node*0.8, 'f', -1, 64) + "m"

	qos_mem_high := strconv.FormatFloat(mem_limit_per_node*0.8, 'f', -1, 64) + "Ki"

	qos_cpu_middle := strconv.FormatFloat(cpu_limit_per_node*0.5, 'f', -1, 64) + "m"

	qos_mem_middle := strconv.FormatFloat(mem_limit_per_node*0.5, 'f', -1, 64) + "Ki"

	cpu_limits := qos_cpu_high

	mem_limits := qos_mem_high

	cpu_requests := qos_cpu_middle

	mem_requests := qos_mem_middle

	var out_qos []byte

	file_str := string(file_byte)

	file_str_list := strings.Split(file_str, "---\n")

	for _, content := range file_str_list {

		if content == "\n" || content == "" {
			continue
		}

		yaml_if := make(map[interface{}]interface{})

		c_byte := []byte(content)

		if err != nil {
			return "", fmt.Errorf("failed to create hpa src: %s", err.Error())
		}

		err = goya.Unmarshal(c_byte, &yaml_if)

		if err != nil {
			return "", fmt.Errorf("failed to create hpa src: %s", err.Error())
		}

		if yaml_if["kind"] == resource_key && yaml_if["metadata"].(map[string]interface{})["name"] == resourcenm {

			src_found = 1

			c_count := len(yaml_if["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{}))

			for j := 0; j < c_count; j++ {

				rsc := map[string]map[string]string{
					"limits": {
						"cpu":    cpu_limits,
						"memory": mem_limits,
					},
					"requests": {
						"cpu":    cpu_requests,
						"memory": mem_requests,
					},
				}

				yaml_if["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[j].(map[string]interface{})["resources"] = rsc

			}

			b_yaml_if, err := goya.Marshal(yaml_if)

			if err != nil {
				return "", fmt.Errorf("failed to create qos src: %s", err.Error())
			}

			out_qos = b_yaml_if

			break
		}

	}

	if src_found == 0 {
		return "", fmt.Errorf("failed to create qos src: %s", "matching key not found")
	}

	err = os.WriteFile("done_question_mark_qos.yaml", out_qos, 0644)

	if err != nil {
		return "", fmt.Errorf("failed to create qos src: %s", err.Error())
	}

	return ".usr/qos_src.yaml", nil

}

func main() {

	// callApiDefStructure()

	// sliceTest()

	// komposeTest()

	//dockercomposeyamlTest()

	// delresourceTest()

	/*
		if s, e := hpaTest(); e != nil {

			fmt.Println(e.Error())

		} else {
			fmt.Println(string(s))
		}
	*/
	if s, e := qosTest(); e != nil {

		fmt.Println(e.Error())

	} else {
		fmt.Println(string(s))
	}
}
