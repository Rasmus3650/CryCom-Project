import numpy as np
import time
from dataclasses import dataclass
from rich.console import Console
from rich.table import Table
from rich.panel import Panel
from rich import box


@dataclass
class GSWParams:
    """Parameters for GSW homomorphic encryption scheme"""
    def __init__(self, q, n, m = None, error_bound=3, t_dist=3):
        self.q = q
        self.n = n
        self.k = n+1
        if m:
            self.m = m
        else:
            self.m = int(2*self.n *np.log2(q)) +1
        self.error_bound = error_bound
        self.t_dist = t_dist
        self.ell = int(np.log2(self.q)) + 1
        self.N = self.k * self.ell

class GSWScheme:
    def __init__(self, params: GSWParams):
        self.params = params
        self.public_key = None
        self.secret_key = None
        # Precompute values for bit decomposition
        self.powers = np.array([pow(2, i, self.params.q) for i in range(self.params.ell)])
        
    def powers_of_two(self, b):
        """Calculate powers of two modulo q for given vector b"""
        result = []
        for i in range(self.params.k):
            b_mod = b[i] % self.params.q
            for power in self.powers:
                result.append((power * b_mod) % self.params.q)
        return np.array(result)

    # --- VECTORIZED IMPLEMENTATION (Third optimization, smaller data types) ---
    def bit_decomp(self, a):
        """Ultra-fast bit decomposition"""
        def vec_decomp(a):
            return np.bitwise_and(
                a[:, np.newaxis] >> np.arange(self.params.ell), 1
            ).reshape(-1)
        return vec_decomp(a) if a.ndim == 1 else np.array([vec_decomp(row) for row in a])

    def bit_decomp_inv(self, a):
        """Ultra-fast inverse bit decomposition"""
        def vec_decomp_inv(a):
            return np.sum(a.reshape(self.params.k, self.params.ell) * self.powers, axis=1) % self.params.q
        return vec_decomp_inv(a) if a.ndim == 1 else np.array([vec_decomp_inv(row) for row in a])

    # ---- VECTORIZED IMPLEMENTATION (Second optimization) ----
    
    # def bit_decomp(self, a):
    #     """Ultra-fast bit decomposition"""
    #     # Use numpy's right shift and bitwise AND for vectorized bit extraction
    #     def vec_decomp(a):
    #         # Vectorized bit extraction using bitwise operations
    #         return np.bitwise_and(
    #             a[:, np.newaxis] >> np.arange(self.params.ell), 1
    #         ).reshape(-1)
    #     # Single-line handling of 1D and 2D inputs
    #     return vec_decomp(np.asarray(a)) if a.ndim == 1 else np.array([vec_decomp(row) for row in a])
    
    # def bit_decomp_inv(self, a):
    #     """Ultra-fast inverse bit decomposition"""
    #     # Vectorized computation using matrix multiplication
    #     def vec_decomp_inv(a):
    #         return np.sum(a.reshape(self.params.k, self.params.ell) * self.powers[:self.params.ell], axis=1) % self.params.q

    #     # Single-line handling of 1D and 2D inputs
    #     return vec_decomp_inv(np.asarray(a)) if a.ndim == 1 else np.array([vec_decomp_inv(row) for row in a])


    # ---- VECTORIZED IMPLEMENTATION (First optimization) ----

    # def bit_decomp(self, a):
    #     """Bit decomposition of input a"""
    #     def vec_decomp(a):
    #         decomp = np.array([
    #             [int(bool(x & (1 << i))) for i in range(self.params.ell)] for x in a
    #         ]).flatten()
    #         return decomp
    #     if type(a) is list:
    #         a = np.array(a)
    #     if a.ndim == 1:  # If a is a vector
    #         return np.array(vec_decomp(a))
    #     elif a.ndim == 2:  # If a is a matrix
    #         return np.array([vec_decomp(row) for row in a])

    # def bit_decomp_inv(self, a):
    #     """Inverse bit decomposition"""
    #     def vec_decomp_inv(a):
    #         result = np.zeros(self.params.k)
    #         # Use numpy vectorized operations instead of explicit loops
    #         chunks = a.reshape(self.params.k, self.params.ell)
    #         # Vectorized computation of chunk sums
    #         result = np.sum(chunks * self.powers, axis=1) % self.params.q
    #         return result
    #     if isinstance(a, list):
    #         a = np.array(a)
    #     if a.ndim == 1:  # vector
    #         return vec_decomp_inv(a)
    #     elif a.ndim == 2:  # matrix
    #         return np.array([vec_decomp_inv(row) for row in a])

    # ---- OLD IMPLEMENTATION ----
    
    def old_bit_decomp(self, a):
        """Bit decomposition of input a"""
        def vec_decomp(a):
            result = []
            for elem in a:
                rev_bit_str = f'{elem%self.params.q:0{self.params.ell}b}'[::-1]
                for digit in rev_bit_str:
                    result.append(int(digit))
            return result
        if type(a) is list:
            a = np.array(a)
        if a.ndim == 1:  # If a is a vector
            return np.array(vec_decomp(a))
        elif a.ndim == 2:  # If a is a matrix
            return np.array([vec_decomp(row) for row in a])
        
    def old_bit_decomp_inv(self, a):
        """Inverse bit decomposition"""
        def vec_decomp_inv(a):
            result = []
            # Precompute powers of 2 mod q to avoid large numbers
            for i in range(self.params.k):
                chunk = a[i * self.params.ell:(i+1) * self.params.ell]
                # Calculate sum using modular arithmetic at each step
                chunk_sum = 0
                for j, val in enumerate(chunk):
                    # Use precomputed power of 2 and multiply mod q
                    term = (self.powers[j] * val) % self.params.q
                    chunk_sum = (chunk_sum + term) % self.params.q
                result.append(chunk_sum)
            return result
        if isinstance(a, list):
            a = np.array(a)
        if a.ndim == 1:  # vector
            return vec_decomp_inv(a)
        elif a.ndim == 2:  # matrix
            return np.array([vec_decomp_inv(row) for row in a])

    def flatten(self, a):
        """Flatten operation"""
        return self.bit_decomp(self.bit_decomp_inv(a))
    
    def old_flatten(self, a):
        """Flatten operation"""
        return np.array(self.old_bit_decomp(self.old_bit_decomp_inv(a)))

    def generate_keys(self):
        """Generate public and secret keys"""
        # Secret Key Gen
        t = np.random.randint(0, self.params.q, size=(self.params.n))
        s = np.array([1] + [-ti for ti in t]) % self.params.q
        v = self.powers_of_two(s)

        # Public Key Gen
        B = np.random.randint(0, self.params.q, size=(self.params.m, self.params.n))
        # Simulated Discrete Gaussian Sampling
        e = np.int64(np.random.normal(-self.params.error_bound,self.params.error_bound,size=(self.params.m)))
        b = ((B@t)+e) % self.params.q
        A = np.column_stack((b, B))
        self.public_key = A
        self.secret_key = v
        return A, v

    # ---- OPTIMIZED IMPLEMENTATION ----
    # def encrypt(self, mu):
    #     R_pk = np.dot(np.random.randint(0, 2, size=(self.params.N, self.params.m), dtype=np.int8), self.public_key)
    #     encrypted_matrix = np.array(self.bit_decomp(R_pk), dtype=np.int64)
    #     np.fill_diagonal(encrypted_matrix, np.diagonal(encrypted_matrix) + mu)
    #     return self.flatten(encrypted_matrix)
        
    def encrypt(self, mu):
        R = np.random.randint(0, 2, size=(self.params.N, self.params.m), dtype=np.uint8)
        # Compute matrix multiplication with smallest possible dtype
        R_pk = (R @ self.public_key) % self.params.q

        # Create sparse identity matrix with minimal memory
        IN = np.eye(self.params.N, dtype=np.uint8)
        # Combine operations to reduce intermediate allocations
        encrypted_matrix = np.array((mu * IN + self.bit_decomp(R_pk))) % self.params.q
        return self.flatten(encrypted_matrix)

    # Old Encrypt

    def old_encrypt(self, mu):
        """Encrypt a message mu"""
        if self.public_key is None:
            raise ValueError("Keys must be generated before encryption")
        # Generate a random binary matrix R
        R = np.random.randint(0, 2, size=(self.params.N, self.params.m))

        # Create a sparse identity matrix
        IN = np.eye(self.params.N, dtype=np.bool)       
        # Perform the encryption
        encrypted_matrix = mu * IN + self.old_bit_decomp((R @ self.public_key) % self.params.q) % self.params.q
        return self.old_flatten(encrypted_matrix)

    def decrypt(self, C):
        """Decrypt ciphertext C"""
        if self.secret_key is None:
            raise ValueError("Keys must be generated before decryption")
        total_error = 0
        for i, vi in enumerate(self.secret_key):
            if self.params.q / 4 < vi <= self.params.q / 2:
                temp_val = round((C[i]@self.secret_key)/vi)
                ref_val = (C[i]@self.secret_key)/vi # Only used to measure the error for experiments
                total_error += abs(temp_val - ref_val)
                return temp_val % 2, total_error

    def mp_decrypt(self, C):
        w = np.zeros(self.params.ell -1)
        for i, row in enumerate(C[:self.params.ell-1]):
            temp = np.int64(0)
            for j, elem in enumerate(row):
                temp += (elem * self.secret_key[j])
            w[i] = temp % self.params.q
        result, counter = [], 0
        for i in range(self.params.ell-2, -1, -1):
            tempx = round(w[i] / self.powers[i])
            if tempx < 1000:    # Error in mp_decrypt, we loop around the ring q, and end up getting invalid decrypt tempx values become much larger than they should be 
                                # but for the wrong decryption we get values 1000+
                x = ((tempx & (2**counter)) >> counter)
                result.append(x)
                for idx in range(i):
                    w[idx] -= (x*(2**idx+counter))
                    w[idx] = w[idx] % self.params.q
            else:
                result.append(0)
            counter += 1
        return int(''.join(map(str, result[::-1])), 2)


    
    def AND(self, C1, C2):
        return self.flatten((C1 @ C2) % self.params.q)

    def XOR(self,C1, C2):
        return self.flatten(C1 + C2)

    def NAND(self, C1,C2):
        return self.flatten((np.eye(C1.shape[0], C2.shape[1], dtype=np.int8)-C1@C2) % self.params.q)

    def MULT_Const(self, C, alpha):
        return self.flatten(self.flatten((alpha * np.eye(C.shape[0], C.shape[1], dtype=np.int8))@C) % self.params.q)


def nice_main():
    console = Console()

    for i in range(1,11):
        console.print(f"\n[bold cyan]====== Iteration {i+1} ======[/bold cyan]")
        
        # Initialize parameters
        start = time.time()
        params = GSWParams(
            q=134217728,
            n=2**(i),
            m=2**(i),
        )
        param_time = time.time() - start

        # Parameter table
        param_table = Table(title="Parameters", box=box.ROUNDED)
        param_table.add_column("Parameter", style="cyan")
        param_table.add_column("Value", style="green")
        param_table.add_row("q", f"{params.q:_.0f}".replace('_', '.'))
        param_table.add_row("n", f"{params.n:_.0f}".replace('_', '.'))
        param_table.add_row("m", f"{params.m:_.0f}".replace('_', '.'))
        console.print(param_table)

        # Create timing table
        timing_table = Table(title="Operation Timings", box=box.ROUNDED)
        timing_table.add_column("Operation", style="cyan")
        timing_table.add_column("Time (seconds)", style="green")
        timing_table.add_row("Parameter Initialization", f"{param_time:.6f}")

        # Create scheme instance
        start = time.time()
        gsw = GSWScheme(params)
        scheme_time = time.time() - start
        timing_table.add_row("Scheme Creation", f"{scheme_time:.6f}")
        
        # Generate keys
        start = time.time()
        _,_ = gsw.generate_keys()
        key_time = time.time() - start
        timing_table.add_row("Key Generation", f"{key_time:.6f}")
        
        # Encrypt a message
        message = np.random.randint(0, 2, dtype=np.int8)
        start = time.time()
        ciphertext = gsw.encrypt(message)
        encrypt_time = time.time() - start
        timing_table.add_row("Encryption", f"{encrypt_time:.6f}")
        
        # Decrypt the message
        start = time.time()
        decrypted, _ = gsw.decrypt(ciphertext)
        decrypt_time = time.time() - start
        timing_table.add_row("Decryption", f"{decrypt_time:.6f}")
        
        console.print(timing_table)

        # Results panel
        results = [
            f"Original message: {message}",
            f"Decrypted message: {decrypted}",
            f"Success: {'[PASS]' if message == decrypted else '[FAIL]'}"
        ]
        # console.print(Panel(
        #     "\n".join(results),
        #     title="Results",
        #     border_style="green" if message == decrypted else "red"
        # )) 
        # Calculate width based on content length or set a fixed width
        max_width = min(max(len(line) for line in results) + 10, 50)  # 150 is an example max width
        console.print(Panel(
            "\n".join(results),
            title="Results",
            border_style="green" if message == decrypted else "red",
            width=max_width
        ))

# Example usage:
def main():
    # Initialize parameters
    start = time.time()
    params = GSWParams(
        q=134217728,
        n=512,
        m=512,
    )
    end = time.time()
    print(f"Parameter initialization time: {end - start:.6f} seconds")
    
    # Create scheme instance
    start = time.time()
    gsw = GSWScheme(params)
    end = time.time()
    print(f"Scheme instance creation time: {end - start:.6f} seconds")
    
    # Generate keys
    start = time.time()
    _,_ = gsw.generate_keys()
    end = time.time()
    print(f"Key generation time: {end - start:.6f} seconds")
    
    # Encrypt a message
    message = 1
    start = time.time()
    ciphertext = gsw.encrypt(message)
    end = time.time()
    print(f"Encryption time: {end - start:.6f} seconds")
    
    # Decrypt the message
    start = time.time()
    decrypted = gsw.decrypt(ciphertext)
    end = time.time()
    print(f"Decryption time: {end - start:.6f} seconds")
    
    # Print results
    print(f"Original message: {message}")
    print(f"Decrypted message: {decrypted}")
    print(f"Successful decryption: {message == decrypted}")

if __name__ == "__main__":
    nice_main()
    #main()
    #test()
