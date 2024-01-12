# Introduction
This repo provides a virtual environment, functioning as a *digital twin* of a cloud-edge infrastructure, specifically designed to enable AI agents to automate and optimize IT operations. The environment represents a *virtual world* where AI agents can both operate and learn. 

The virtual environment serves multiple purposes:

* **Data collection:** Connect to a real infrastructure to gather real-time data, essential for training AI agents.
* **Simulation:** Utilizing the collected data, AI agents can be trained within the virtual environment, allowing them to learn and adapt to various situations without impacting the real infrastructure.
* **Deployment:** Post-training, by connecting to a real infrastrucure, the AI agents can be deployed to leverage their learned strategies and to automate and enhance IT operations.

The environment can be used to generate *state and action spaces*. 
* The state-space represents configurations (e.g. hosts, VMs, CPU load etc) that an AI agent might encounter in the environment. Each state is a unique snapshot of the environment at a given time. 
* The action space defines the set of all possible actions (e.g. place or scale a VM) that the AI agent can take at any given state.  
