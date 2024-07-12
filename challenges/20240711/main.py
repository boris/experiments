import subprocess
import json

def running_containers():
    containers = subprocess.run(["docker", "ps", "--format", "{{json .}}"], capture_output=True, text=True)
    json_lines = containers.stdout.strip().split('\n')


    for l in json_lines:
        container_info = json.loads(l)
        container_id = container_info['ID']

        stats_result = subprocess.run(["docker", "stats", "--no-stream", "--format", "{{json .}}", container_id], capture_output=True, text=True)

        container_stats = json.loads(stats_result.stdout)

        print(f"{container_id:<12} {container_info['Names']:<20} {container_info['Status']:<20} {container_stats['CPUPerc']:<10} {container_stats['MemUsage']}")



if __name__ == "__main__":
    running_containers()
