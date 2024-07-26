import re
from datetime import datetime
import matplotlib.pyplot as plt
import matplotlib.dates as mdates

def parse_log_file(file_path):
    data = {
        'resource_utilization': [],
        'qos_violations': [],
        'container_migrations': []
    }
    
    with open(file_path, 'r') as f:
        for line in f:
            timestamp_match = re.search(r'(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})', line)
            if timestamp_match:
                timestamp = datetime.strptime(timestamp_match.group(1), '%Y/%m/%d %H:%M:%S')
                
                # Parse resource utilization
                if 'CPU:' in line and 'Memory:' in line and 'GPU:' in line:
                    try:
                        cpu = float(re.search(r'CPU: ([\d.]+)%', line).group(1))
                        memory = float(re.search(r'Memory: ([\d.]+)%', line).group(1))
                        gpu = float(re.search(r'GPU: ([\d.]+)%', line).group(1))
                        data['resource_utilization'].append((timestamp, cpu, memory, gpu))
                    except (AttributeError, ValueError):
                        print(f"Error parsing resource utilization on line: {line}")
                
                # Parse QoS violations
                elif 'QoS Violations:' in line:
                    try:
                        violations = int(re.search(r'QoS Violations: (\d+)', line).group(1))
                        data['qos_violations'].append((timestamp, violations))
                    except (AttributeError, ValueError):
                        print(f"Error parsing QoS violations on line: {line}")
                
                # Parse container migrations
                elif 'Migrated container' in line:
                    try:
                        container_id = re.search(r'Migrated container (\S+)', line).group(1)
                        source_host = re.search(r'from host (\S+)', line).group(1)
                        dest_host = re.search(r'to host (\S+)', line).group(1)
                        data['container_migrations'].append((timestamp, container_id, source_host, dest_host))
                    except (AttributeError, ValueError):
                        print(f"Error parsing container migrations on line: {line}")
    
    return data

def create_resource_utilization_graph(data, strategy_name):
    if not data['resource_utilization']:
        print(f"No resource utilization data available for {strategy_name}.")
        return
    
    timestamps, cpu, memory, gpu = zip(*data['resource_utilization'])
    
    plt.figure(figsize=(12, 6))
    plt.plot(timestamps, cpu, label='CPU')
    plt.plot(timestamps, memory, label='Memory')
    plt.plot(timestamps, gpu, label='GPU')
    
    plt.xlabel('Time')
    plt.ylabel('Utilization (%)')
    plt.title(f'Resource Utilization Over Time - {strategy_name}')
    plt.legend()
    
    plt.gca().xaxis.set_major_formatter(mdates.DateFormatter('%H:%M:%S'))
    plt.gcf().autofmt_xdate()
    
    plt.tight_layout()
    plt.savefig(f'{strategy_name}_resource_utilization.png')
    plt.close()

def create_qos_violations_graph(data, strategy_name):
    if not data['qos_violations']:
        print(f"No QoS violations data available for {strategy_name}.")
        return
    
    timestamps, violations = zip(*data['qos_violations'])
    
    plt.figure(figsize=(12, 6))
    plt.bar(timestamps, violations, width=0.01)
    
    plt.xlabel('Time')
    plt.ylabel('Number of QoS Violations')
    plt.title(f'QoS Violations Over Time - {strategy_name}')
    
    plt.gca().xaxis.set_major_formatter(mdates.DateFormatter('%H:%M:%S'))
    plt.gcf().autofmt_xdate()
    
    plt.tight_layout()
    plt.savefig(f'{strategy_name}_qos_violations.png')
    plt.close()

def create_container_migrations_graph(data, strategy_name):
    if not data['container_migrations']:
        print(f"No container migrations data available for {strategy_name}.")
        return
    
    timestamps = [t for t, _, _, _ in data['container_migrations']]
    migrations = range(1, len(timestamps) + 1)
    
    plt.figure(figsize=(12, 6))
    plt.plot(timestamps, migrations, marker='o')
    
    plt.xlabel('Time')
    plt.ylabel('Cumulative Number of Migrations')
    plt.title(f'Container Migrations Over Time - {strategy_name}')
    
    plt.gca().xaxis.set_major_formatter(mdates.DateFormatter('%H:%M:%S'))
    plt.gcf().autofmt_xdate()
    
    plt.tight_layout()
    plt.savefig(f'{strategy_name}_container_migrations.png')
    plt.close()

def process_strategy(strategy_name):
    log_file = f'{strategy_name}_simulation.log'
    data = parse_log_file(log_file)

    print("Data was parsed successfully from logs")
    
    create_resource_utilization_graph(data, strategy_name)
    create_qos_violations_graph(data, strategy_name)
    create_container_migrations_graph(data, strategy_name)

def main():
    strategies = ['RoundRobin', 'Priority', 'BinPacking']
    for strategy in strategies:
        process_strategy(strategy)
    
    print("Visualization complete. Check the current directory for the generated graphs.")

if __name__ == "__main__":
    main()
