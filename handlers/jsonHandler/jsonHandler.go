package jsonHandler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"surface_attack/consts"
	"surface_attack/models"
	"surface_attack/providers/crudProvider"
)

type InputJsonHandler struct {
	CrudProvider crudProvider.Interface
}

func (j *InputJsonHandler) OnBoot(location string) error {

	log.Println("begin bootup json data processing")
	inputData, e := j.populateInputDataObject(location)
	if e != nil {
		return e
	}

	destToSources := j.getDestToSourcesMap(inputData.Fw_rules)

	vmToSourceTags, vmToSourceVmsRaw, tagToVms := j.getMapsFromVmArray(inputData.Vms, destToSources)

	j.populateVmToSourceVmsRaw(tagToVms, vmToSourceVmsRaw, vmToSourceTags)

	vmToSourceVmsFinal := j.getFinalMapFromRaw(vmToSourceVmsRaw)

	log.Println("json data processing succesfully finished")
	for k, v := range vmToSourceVmsFinal {
		e := j.CrudProvider.Set(k, strings.Join(v, ","))
		if e != nil {
			return e
		}
	}

	vmCountInt := len(vmToSourceVmsFinal)
	vmCountStr := strconv.Itoa(vmCountInt)

	e = j.CrudProvider.Set(consts.REDIS_VM_COUNT_KEY, vmCountStr)
	if e != nil {
		return e
	}
	log.Println("succesfully finished uploading json data to redis")
	return nil
}

//read input json document into memory object inputData
func (j *InputJsonHandler) populateInputDataObject(location string) (*models.InputJson, error) {

	plan, _ := ioutil.ReadFile(location)
	var inputData models.InputJson
	err := json.Unmarshal(plan, &inputData)
	if err != nil {
		return nil, errors.New("InputJsonHandler -> populateInputDataObject: " + err.Error())
	}

	return &inputData, nil
}

//for each fw_rule create map of "dest-tag"->["source tag"...]
func (j *InputJsonHandler) getDestToSourcesMap(fwRules []models.Fw_rules) map[string][]string {
	var destToSources = make(map[string][]string)

	for _, fw_rule := range fwRules {
		if _, ok := destToSources[fw_rule.Dest_tag]; !ok {
			destToSources[fw_rule.Dest_tag] = []string{}
		}

		if j.alreadyExistInArray(destToSources[fw_rule.Dest_tag], fw_rule.Source_tag) {
			continue
		}
		destToSources[fw_rule.Dest_tag] = append(destToSources[fw_rule.Dest_tag], fw_rule.Source_tag)
	}
	return destToSources
}

//for each vm create maps of "vm"->["source tag"...] && "tag"->["vm"...].
// then use these two maps to create a map of "vm" ->["source vm"...]
//initialize all maps used
func (j *InputJsonHandler) getMapsFromVmArray(vms []models.Vm, destToSources map[string][]string) (vmToSourceTags, vmToSourceVmsRaw, tagToVms map[string][]string) {

	//initialize all maps used
	tagToVms = make(map[string][]string)
	vmToSourceTags = make(map[string][]string)
	vmToSourceVmsRaw = make(map[string][]string)

	for _, vm := range vms {

		vmToSourceTags[vm.Vm_id] = []string{}
		vmToSourceVmsRaw[vm.Vm_id] = []string{}

		for _, tag := range vm.Tags {

			for _, source := range destToSources[tag] {
				if j.alreadyExistInArray(vmToSourceTags[vm.Vm_id], source) {
					continue
				}
				vmToSourceTags[vm.Vm_id] = append(vmToSourceTags[vm.Vm_id], source)
			}

			if _, ok := tagToVms[tag]; !ok {
				tagToVms[tag] = []string{}
			}
			//TAG -> [vms...]
			tagToVms[tag] = append(tagToVms[tag], vm.Vm_id)
		}
	}
	return vmToSourceTags, vmToSourceVmsRaw, tagToVms
}

//create final map with key vm and value array of all possible source vms
func (j *InputJsonHandler) populateVmToSourceVmsRaw(tagToVms, vmToSourceVmsRaw, vmToSourceTags map[string][]string) {
	for vm, sourceTags := range vmToSourceTags {
		for _, sourceTag := range sourceTags {
			vmToSourceVmsRaw[vm] = append(vmToSourceVmsRaw[vm], tagToVms[sourceTag]...)
		}
	}
}

func (j *InputJsonHandler) getFinalMapFromRaw(vmToSourceVmsRaw map[string][]string) (vmToSourceVmsFinal map[string][]string) {
	vmToSourceVmsFinal = make(map[string][]string)

	for vm, sourceVms := range vmToSourceVmsRaw {
		vmToSourceVmsFinal[vm] = []string{}
		for _, sourceVm := range sourceVms {
			if vm != sourceVm {
				if j.alreadyExistInArray(vmToSourceVmsFinal[vm], sourceVm) {
					continue
				}
				vmToSourceVmsFinal[vm] = append(vmToSourceVmsFinal[vm], sourceVm)
			}
		}
	}
	return vmToSourceVmsFinal
}

func (j *InputJsonHandler) alreadyExistInArray(arr []string, tag string) bool {

	for _, v := range arr {
		if v == tag {
			return true
		}
	}
	return false
}
