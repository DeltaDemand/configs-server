package model

type Group map[string]ConfigItems  //根据[AgentName]找到对应Agent拥有ConfigItems
type ConfigItems map[string]string //根据每个配置项的名字，如cpu_rate找到对应配置value
