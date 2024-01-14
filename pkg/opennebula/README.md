The OpenNebula scheduler will send this request to the AI Orchestrator when a VM is going to be deployed.

```json
{
  "VMS": [
    {
      "CAPACITY": {
        "CPU": 1.0,
        "DISK_SIZE": 2252,
        "MEMORY": 786432
      },
      "HOST_IDS": [
        0,
        2,
        3,
        4
      ],
      "ID": 7,
      "STATE": "PENDING",
      "USER_TEMPLATE": {
        "LOGO": "images/logos/ubuntu.png",
        "LXD_SECURITY_PRIVILEGED": "true",
        "SCHED_REQUIREMENTS": "ID=\"0\" | ID=\"2\" | ID=\"3\" | ID=\"4\""
      }
    }
  ]
}
```

The server is expecting this JSON response:

```json
{
  "VMS": [
    {
      "ID": 7,
      "HOST_ID": 4
    }
  ]
}
```
