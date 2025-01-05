import unittest
import numpy as np
import math
from gsw_scheme import GSWParams, GSWScheme

class TestGSWHomomorphicEncryption(unittest.TestCase):
    def setUp(self):
        """Set up test parameters that will be used across multiple tests"""
        self.params = GSWParams(
            q=134217728,
            n=64,
            m=64,
        )
        self.gsw = GSWScheme(self.params)
        self.gsw.generate_keys()

    def test_aux_properties(self):
        """Test properties of the auxiliary functions"""
        a = np.random.randint(0, self.params.q, size=self.params.k)
        b = np.random.randint(0, self.params.q, size=self.params.k)
        self.assertEqual(self.gsw.bit_decomp(a)@self.gsw.powers_of_two(b) % self.params.q,(a@b)%self.params.q)
        self.assertEqual(self.gsw.bit_decomp(a)@self.gsw.powers_of_two(b) % self.params.q,(self.gsw.flatten(self.gsw.bit_decomp(a))@self.gsw.powers_of_two(b))%self.params.q)
        self.assertEqual((self.gsw.bit_decomp(a)@self.gsw.powers_of_two(b)) % self.params.q,(self.gsw.flatten(self.gsw.bit_decomp(a))@self.gsw.powers_of_two(b))%self.params.q)
        self.assertTrue(
            np.array_equal(
                self.gsw.public_key,
                self.gsw.bit_decomp_inv(self.gsw.bit_decomp(self.gsw.public_key))
            )
        )    
    def test_basic_encryption_decryption(self):
        """Test basic encryption and decryption functionality"""
        message = 1
        ciphertext = self.gsw.encrypt(message)
        decrypted = self.gsw.decrypt(ciphertext)
        self.assertEqual(message, decrypted, "Encryption/decryption failed for message=1")

    def test_homomorphic_multiplication(self):
        """Test homomorphic multiplication (AND operation)"""
        # Encrypt inputs

        c1 = self.gsw.encrypt(1)
        c2 = self.gsw.encrypt(0)
        
        # Test 0 XOR 0 = 0
        c3 = self.gsw.XOR(c2, c2)
        self.assertEqual(0, self.gsw.decrypt(c3), "Operation failed for 0 XOR 0")
        
        # Test 1 XOR 1 = 1
        c4 = self.gsw.XOR(c1, c1)
        self.assertEqual(0, self.gsw.decrypt(c4), "Operation failed for 1 XOR 1")

        # Test 1 XOR 0 = 0
        c4 = self.gsw.XOR(c1, c2)
        self.assertEqual(1, self.gsw.decrypt(c4), "Operation failed for 1 XOR 0")

    def test_homomorphic_addition(self):
        """Test homomorphic addition (XOR operation)"""
        # Encrypt inputs
        c1 = self.gsw.encrypt(1)
        c2 = self.gsw.encrypt(0)
        
        # Test 1 AND 0 = 0
        c3 = self.gsw.AND(c1, c2)
        self.assertEqual(0, self.gsw.decrypt(c3), "Operation failed for 1 AND 0")
        
        # Test 1 AND 1 = 1
        c4 = self.gsw.AND(c1, c1) 
        self.assertEqual(1, self.gsw.decrypt(c4), "Operation failed for 1 AND 1")

        # Test 0 AND 0 = 0
        c5 = self.gsw.AND(c2, c2)
        self.assertEqual(0, self.gsw.decrypt(c5), "Operation failed for 0 AND 0")

    def test_bit_decomposition_properties(self):
        """Test properties of bit decomposition operations"""
        k = self.params.n + 1
        # Generate random vectors for testing
        a = np.random.randint(0, self.params.q, size=k)
        b = np.random.randint(0, self.params.q, size=k)
        
        # Test first property: bit_decomp(a) @ powers_of_two(b) ≡ a @ b (mod q)
        a_decomp = self.gsw.bit_decomp(a)
        b_powers = self.gsw.powers_of_two(b)
        left_side = (a_decomp @ b_powers) % self.params.q
        right_side = (a @ b) % self.params.q
        self.assertTrue(np.array_equal(left_side, right_side), 
                       "First bit decomposition property failed")
        
        # Test second property: bit_decomp_inv(bit_decomp(a)) = a
        reconstructed = self.gsw.bit_decomp_inv(self.gsw.bit_decomp(a))
        self.assertTrue(np.array_equal(a % self.params.q, reconstructed), 
                       "Bit decomposition inverse property failed")

def main():
    # Run the tests
    unittest.main()
if __name__ == '__main__':
    main()
