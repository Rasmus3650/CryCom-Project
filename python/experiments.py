import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import gsw_scheme as gsw
import time
from memory_profiler import memory_usage

def measure_avg_runtime(func, iterations=1):
    """Helper to measure average runtime of a function."""
    times = []
    for _ in range(iterations):
        start_time = time.time()
        func()
        times.append(time.time() - start_time)
    return sum(times) / len(times)


def collect_and_plot_comparison():
    # Initialize parameters
    q = 134217728
    dimensions = []
    
    # Initialize timing dictionaries
    avg_new_times = {"key_gen": [], "encryption": [], "decryption": []}
    avg_old_times = {"key_gen": [], "encryption": [], "decryption": []}
    
    # Collect Python implementation data
    for ctr in range(1, 9):
        n = 2**ctr
        m = n
        #m = 2*n*np.log2(q) + 1
        params = gsw.GSWParams(q=q, n=n, m=m)
        plaintext = np.random.randint(0, 2)
        print(n)
        # Initialize schemes
        gsw_new = gsw.GSWScheme(params)
        gsw_old = gsw.GSWScheme(params)

        def new_encrypt():
            new_encrypt.ctext = gsw_new.encrypt(plaintext)
            return new_encrypt.ctext

        # Measure average runtimes for new scheme
        avg_new_times["key_gen"].append(measure_avg_runtime(lambda: gsw_new.generate_keys()))
        avg_new_times["encryption"].append(measure_avg_runtime(new_encrypt))
        ctext = new_encrypt.ctext
        avg_new_times["decryption"].append(measure_avg_runtime(lambda: gsw_new.decrypt(ctext)))
        print(f"{n} opt done")
        def old_encrypt():
            old_encrypt.ctext = gsw_old.old_encrypt(plaintext)
            return old_encrypt.ctext

        # Measure average runtimes for old scheme
        avg_old_times["key_gen"].append(measure_avg_runtime(lambda: gsw_old.generate_keys()))
        avg_old_times["encryption"].append(measure_avg_runtime(old_encrypt))
        temp_ctext = old_encrypt.ctext
        avg_old_times["decryption"].append(measure_avg_runtime(lambda: gsw_old.decrypt(temp_ctext)))
        print(f"{n} vanilla done")
        dimensions.append(n)
    
    # Calculate total times for Python implementations
    python_new_total = [sum(x) for x in zip(avg_new_times["key_gen"], 
                                          avg_new_times["encryption"], 
                                          avg_new_times["decryption"])]
    python_old_total = [sum(x) for x in zip(avg_old_times["key_gen"], 
                                         avg_old_times["encryption"], 
                                         avg_old_times["decryption"])]

    # Load and process Go data
    go_vanilla_data = pd.read_csv("experiment_data/go_vanilla.csv")
    go_optimized_data = pd.read_csv("experiment_data/go_optimized.csv")
    
    # Convert Go times from ns to seconds
    go_vanilla_data['time_in_s'] = go_vanilla_data['time_in_ns'] / 1e9
    go_optimized_data['time_in_s'] = go_optimized_data['time_in_ns'] / 1e9

    # Create the visualization
    plt.figure(figsize=(12, 8))
    
    # Plot all implementations
    plt.plot(dimensions, python_new_total, 'o-', label='Python (Optimized)', color='blue', linewidth=2)
    plt.plot(dimensions, python_old_total, 's-', label='Python (Vanilla)', color='red', linewidth=2)
    
    # Filter Go data for the same dimensions as Python
    go_vanilla_filtered = go_vanilla_data[go_vanilla_data['n'].isin(dimensions)]
    go_optimized_filtered = go_optimized_data[go_optimized_data['n'].isin(dimensions)]
    
    plt.plot(go_vanilla_filtered['n'], go_vanilla_filtered['time_in_s'], '^-', 
             label='Go (Vanilla)', color='green', linewidth=2)
    plt.plot(go_optimized_filtered['n'], go_optimized_filtered['time_in_s'], 'D-', 
             label='Go (Optimized)', color='purple', linewidth=2)

    plt.xlabel('Dimension Size (n)', fontsize=12)
    plt.ylabel('Time (seconds)', fontsize=12)
    plt.title('Performance Comparison: Python vs Go Implementations', fontsize=14)
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.legend(fontsize=10)
    
    # Set log scales but with custom tick labels
    plt.yscale('log')
    plt.xscale('log', base=2)
    
    # Customize x-axis ticks to show actual n values
    plt.xticks(dimensions, dimensions)
    
    # Customize y-axis to show actual seconds with more decimal places
    yticks = plt.yticks()[0]
    plt.yticks(yticks, [f'{x:.6f}' for x in yticks])
    
    plt.tight_layout()
    return plt


def memory_comparison():
    # Initialize parameters
    q = 134217728
    dimensions = []
    mem_new = []
    mem_old = []
    
    def run_new(params):
        gsw_scheme = gsw.GSWScheme(params)
        gsw_scheme.generate_keys()
        mu = np.random.randint(0, 2)
        ct = gsw_scheme.encrypt(mu)
        ptext = gsw_scheme.decrypt(ct)
        return ptext
    
    def run_old(params):
        old_gsw = gsw.GSWScheme(params)
        old_gsw.generate_keys()
        mu = np.random.randint(0, 2)
        ct = old_gsw.old_encrypt(mu)
        ptext = old_gsw.decrypt(ct)
        return ptext
    
    # Collect Python implementation data
    for ctr in range(1, 7):
        params = gsw.GSWParams(
            q=q,
            n=2**ctr,
            m=2**ctr
            #m = 2*(2**ctr)*np.log2(q) + 1
        )
        # Measure memory usage for the new scheme
        new_mem = memory_usage((run_new, (params,)), max_iterations=1)
        mem_new.append(max(new_mem))
        
        # Measure memory usage for the old scheme
        old_mem = memory_usage((run_old, (params,)), max_iterations=1)
        mem_old.append(max(old_mem))
        dimensions.append(2**ctr)
    
    # Load and process Go data
    go_vanilla_data = pd.read_csv("experiment_data/go_vanilla.csv")
    go_optimized_data = pd.read_csv("experiment_data/go_optimized.csv")
    
    # Convert bytes to MB for Go data
    go_vanilla_data['mem_in_mb'] = go_vanilla_data['total_mem_in_byte'] / (1024 * 1024)
    go_optimized_data['mem_in_mb'] = go_optimized_data['total_mem_in_byte'] / (1024 * 1024)
    
    # Create the visualization
    plt.figure(figsize=(12, 8))
    
    # Plot all implementations
    plt.plot(dimensions, mem_new, 'o-', label='Python (Optimized)', color='blue', linewidth=2)
    plt.plot(dimensions, mem_old, 's-', label='Python (Vanilla)', color='red', linewidth=2)
    
    # Filter Go data for the same dimensions as Python
    go_vanilla_filtered = go_vanilla_data[go_vanilla_data['n'].isin(dimensions)]
    go_optimized_filtered = go_optimized_data[go_optimized_data['n'].isin(dimensions)]
    
    plt.plot(go_vanilla_filtered['n'], go_vanilla_filtered['mem_in_mb'], '^-', 
             label='Go (Vanilla)', color='green', linewidth=2)
    plt.plot(go_optimized_filtered['n'], go_optimized_filtered['mem_in_mb'], 'D-', 
             label='Go (Optimized)', color='purple', linewidth=2)

    plt.xlabel('Dimension Size (n)', fontsize=12)
    plt.ylabel('Memory Usage (MB)', fontsize=12)
    plt.title('Memory Usage Comparison: Python vs Go Implementations', fontsize=14)
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.legend(fontsize=10)
    
    # Set log scales but with custom tick labels
    plt.yscale('log')
    plt.xscale('log', base=2)
    
    # Customize x-axis ticks to show actual n values
    plt.xticks(dimensions, dimensions)
    
    # Customize y-axis to show actual MB values with more precision
    yticks = plt.yticks()[0]
    plt.yticks(yticks, [f'{x:.6f}' for x in yticks])
    
    plt.tight_layout()
    return plt

def error_growth():
    # Initialize parameters
    q = 134217728
    n = 32
    # m = 2 * n * np.log2(q) + 1
    m = n
    params = gsw.GSWParams(q=q, n=n, m=m)
    
    # Initialize schemes
    gsw_new = gsw.GSWScheme(params)    
    # Generate keys
    gsw_new.generate_keys()
    num_wires = [2**i for i in range(1, 8)]
    errors = []
    runtime = []
    validator = []
    for wires in num_wires:
        print(f"Number of wires: {wires}")
        and_binary_errors = []
        and_binary_runtime = []
        and_stack_errors = []
        and_stack_runtime = []

        xor_binary_errors = []
        xor_binary_runtime = []
        xor_stack_errors = []
        xor_stack_runtime = []

        for _ in range(2):
            # Generate plaintexts and encrypt them to get ciphertexts
            ptexts = [np.random.randint(0, 2) for _ in range(wires)]
            print("Plaintexts:", ptexts)
            ctexts = [gsw_new.encrypt(ptext) for ptext in ptexts]
            print("Done encrypting")

            for gate in ["XOR", "AND"]:
                print(f"Processing {gate}")
                # ---------------- Binary Tree Reduction ----------------
                start_time = time.time()
                current_layer = iter(ctexts)  # Use an iterator for efficient processing
                next_layer = []
                while True:
                    try:
                        # Take two ciphertexts at a time from the iterator
                        left_ctext = next(current_layer)
                        right_ctext = next(current_layer)
                        # Perform the gate operation
                        if gate == "XOR":
                            result = gsw_new.XOR(left_ctext, right_ctext)
                        elif gate == "AND":
                            result = gsw_new.AND(left_ctext, right_ctext)
                        next_layer.append(result)
                    except StopIteration:
                        # Move to the next layer
                        if len(next_layer) == 1:
                            # Final root node
                            break
                        current_layer = iter(next_layer)
                        next_layer = []
                binary_runtime = time.time() - start_time
                binary_output, binary_error = gsw_new.decrypt(next_layer[0])
                print(f"Binary runtime: {binary_runtime}")
                print(f"Binary output: {binary_output} - Error: {binary_error}")

                # ---------------- Stack-Based Reduction ----------------
                start_time = time.time()
                stack = ctexts[:]
                while len(stack) > 1:
                    left_ctext = stack.pop()
                    right_ctext = stack.pop()
                    # Perform the gate operation
                    if gate == "XOR":
                        result = gsw_new.XOR(left_ctext, right_ctext)
                    elif gate == "AND":
                        result = gsw_new.AND(left_ctext, right_ctext)
                    stack.append(result)
                stack_runtime = time.time() - start_time
                stack_output, stack_error = gsw_new.decrypt(stack[0])
                print(f"Stack runtime: {stack_runtime}")
                print(f"Stack output: {stack_output} - Error: {stack_error}")
                validator.append(binary_output == stack_output)
                if gate == "AND":
                    and_binary_errors.append(binary_error)
                    and_binary_runtime.append(binary_runtime)
                    and_stack_errors.append(stack_error)
                    and_stack_runtime.append(stack_runtime)
                elif gate == "XOR":
                    xor_binary_errors.append(binary_error)
                    xor_binary_runtime.append(binary_runtime)
                    xor_stack_errors.append(stack_error)
                    xor_stack_runtime.append(stack_runtime)
        vals = {"AND":[sum(and_binary_errors) / len(and_binary_errors) if and_binary_errors else 0, sum(and_stack_errors) / len(and_stack_errors) if and_stack_errors else 0],
                "XOR":[sum(xor_binary_errors) / len(xor_binary_errors) if xor_binary_errors else 0, sum(xor_stack_errors) / len(xor_stack_errors) if xor_stack_errors else 0]}
        errors.append(vals)
        runs = {"AND":[sum(and_binary_runtime) / len(and_binary_runtime) if and_binary_runtime else 0, sum(and_stack_runtime) / len(and_stack_runtime) if and_stack_runtime else 0],
                "XOR":[sum(xor_binary_runtime) / len(xor_binary_runtime) if xor_binary_runtime else 0, sum(xor_stack_runtime) / len(xor_stack_runtime) if xor_stack_runtime else 0]}
        runtime.append(runs)
    return errors, runtime, validator

def plot_gate_performance(errors, runtime):
    """
    Plot runtime and error metrics for different gate operations.
    
    Args:
    errors: List of dicts containing error values for each gate type and implementation
    runtime: List of dicts containing runtime values for each gate type and implementation
    """
    # Calculate number of wires - starting at 2^1 for first data point
    num_wires = [2**(i+1) for i in range(len(errors))]
    
    # Create figure with subplots
    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 10))
    
    # Colors and markers for different implementations
    colors = {'AND': 'blue', 'XOR': 'red'}
    styles = {'binary': '-o', 'stack': '--s'}
    
    # Plot runtime
    for gate in ['AND', 'XOR']:
        binary_runtime = [run[gate][0] for run in runtime]  # Index 0 for binary
        stack_runtime = [run[gate][1] for run in runtime]   # Index 1 for stack
        
        ax1.plot(num_wires, binary_runtime, styles['binary'], color=colors[gate], 
                label=f'{gate} Binary')
        ax1.plot(num_wires, stack_runtime, styles['stack'], color=colors[gate], 
                label=f'{gate} Stack')
    
    ax1.set_title('Runtime Comparison')
    ax1.set_xlabel('Number of Input Wires')
    ax1.set_ylabel('Runtime (seconds)')
    ax1.grid(True)
    ax1.legend()
    
    # Plot errors
    for gate in ['AND', 'XOR']:
        binary_errors = [err[gate][0] for err in errors]    # Index 0 for binary
        stack_errors = [err[gate][1] for err in errors]     # Index 1 for stack
        
        ax2.plot(num_wires, binary_errors, styles['binary'], color=colors[gate], 
                label=f'{gate} Binary')
        ax2.plot(num_wires, stack_errors, styles['stack'], color=colors[gate], 
                label=f'{gate} Stack')
    
    ax2.set_title('Error Comparison')
    ax2.set_xlabel('Number of Input Wires')
    ax2.set_ylabel('Error Rate')
    ax2.set_yscale('log')
    ax2.grid(True)
    ax2.legend()
    
    plt.tight_layout()
    return plt
    

def main():
    #plt = collect_and_plot_comparison()
    #plt.savefig('runtime_graph.png', dpi=300, bbox_inches='tight')
    #plt = memory_comparison()
    #plt.savefig('memory_graph.png', dpi=300, bbox_inches='tight')
    # Run profiling
    
    #results = detailed_memory_profile()


    errors, runtime, validator = error_growth()
    plt = plot_gate_performance(errors, runtime)
    plt.savefig('error_graph.png', dpi=300, bbox_inches='tight')

    #plt.show()
    print("-----------------------------------------------")
    print("Errors:", errors)
    print("Runtimes:", runtime)
    print("Validator:", validator)

main()