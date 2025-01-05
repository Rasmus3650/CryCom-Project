from gsw_scheme import GSWScheme, GSWParams
import numpy as np

counter = 1
class Alice():
    def __init__(self, bits, scheme):
        self.gsw = scheme
        self.A = self.gsw.public_key
        self.v = self.gsw.secret_key
        self.bits = bits
        self.ctexts = [self.gsw.encrypt(bit) for bit in self.bits]
        self.B_ctexts = None

    def share_ctext(self):
        return self.ctexts
    
    def decrypt(self, ctext):
        return self.gsw.decrypt(ctext)

class Bob():
    def __init__(self, bits, scheme):
        self.gsw = scheme
        self.A = self.gsw.public_key
        self.bits = bits
        self.ctexts = [self.gsw.encrypt(bit) for bit in self.bits]
        self.A_ctexts = None

    def share_ctext(self):
        return self.ctexts
    
    def OR(self, x,y):
        return self.gsw.XOR(self.gsw.XOR(x,y),self.gsw.AND(x,y))

    def EQ(self,x,y):
        return self.NOT(self.gsw.XOR(x, y))
    def GT(self, x,y):
        return self.gsw.AND(self.gsw.XOR(x, y), x)
    
    def NOT(self, x):
        val_1 = self.gsw.encrypt(1)
        return self.gsw.XOR(x, val_1)
   
    def eval_circuit(self):
        global counter
        result = self.GT(self.A_ctexts[0], self.ctexts[0])
        for i in range(1, counter):
            eq_chain = self.EQ(self.A_ctexts[0], self.ctexts[0])
            for j in range(1, i):
                eq_chain = self.gsw.AND(eq_chain, self.EQ(self.A_ctexts[j], self.ctexts[j]))
            gt_chain = self.GT(self.A_ctexts[i], self.ctexts[i])
            result = self.OR(result, self.gsw.AND(eq_chain, gt_chain))
        return result

def yao_boolean_verifier(a_bits, b_bits):
    a_num = int("".join(map(str, a_bits)), 2)
    b_num = int("".join(map(str, b_bits)), 2)
    return a_num > b_num

def test_yao():
    global counter    
    total_error = []
    for _ in range(1,5):
        possibilities = []
        for i in range(2**counter):
            binary_representation = [int(bit) for bit in f"{i:0{counter}b}"]
            possibilities.append(binary_representation)
        failed_counter = 0
        for i, recipient in enumerate(possibilities):
            for j, donor in enumerate(possibilities):
                params = GSWParams(q=134217728,n=32, m =32)
                gsw = GSWScheme(params)
                A, v = gsw.generate_keys() # Pretend Bob can only access A
                a = Alice(recipient, gsw)
                b = Bob(donor, gsw)
                b.A_ctexts = a.share_ctext()
                z, error = a.decrypt(b.eval_circuit())
                total_error.append(error)
                
                check = z == yao_boolean_verifier(recipient, donor)
                if not check:
                    failed_counter += 1
                print(z, error, check, recipient, donor)
        print("Average error: ", sum(total_error)/len(total_error))
        print(f"Number of fails for circuit input of {counter} bits: ", failed_counter)
        counter += 1

test_yao()