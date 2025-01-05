from gsw_scheme import GSWScheme, GSWParams
import numpy as np

class Alice():
    def __init__(self, bits, scheme):
        self.gsw = scheme
        self.A = self.gsw.public_key
        self.v = self.gsw.secret_key
        self.bits = [int(bit) for bit in bits]
        #self.bits = [random.randint(0,2) for _ in range(4)]
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
        self.bits = [int(not int(bit)) for bit in bits]
        #self.bits = [random.randint(0,2) for _ in range(4)]
        self.ctexts = [self.gsw.encrypt(bit) for bit in self.bits]
        self.A_ctexts = None

    def share_ctext(self):
        return self.ctexts
    
    def eval_circuit(self):
        lhs = self.gsw.NAND(self.ctexts[0], self.A_ctexts[0])
        mhs = self.gsw.NAND(self.ctexts[1], self.A_ctexts[1])
        rhs = self.gsw.NAND(self.ctexts[2], self.A_ctexts[2])
        return self.gsw.AND(self.gsw.AND(rhs, mhs), lhs)


def convert_type_to_bitstring(blood_type):
    if blood_type not in ["O-", "O+", "B-", "B+", "A-", "A+", "AB-", "AB+"]:
        raise ValueError("Blood type not recognized")
    result_str = ""
    if "AB" in blood_type:
        result_str += "11"
    elif "A" in blood_type:
        result_str += "10"
    elif "B" in blood_type:
        result_str += "01"
    elif "O" in blood_type:
        result_str += "00"
    if "+" in blood_type:
        result_str += "1"
    elif "-" in blood_type:
        result_str += "0"
    return result_str

def blood_compatibility_boolean(recipient_type, donor_type):
    recipient_bits = convert_type_to_bitstring(recipient_type)
    donor_bits = convert_type_to_bitstring(donor_type)
    xa, xb, xr = int(recipient_bits[0]), int(recipient_bits[1]), int(recipient_bits[2])
    ya, yb, yr = int(donor_bits[0]), int(donor_bits[1]), int(donor_bits[2])
    return (((1 ^ ((1 ^ ya) and xa)) and (1 ^ ((1 ^ yb) and xb))) and (1 ^ ((1 ^ yr) and xr)))


def test_blood_type_compatibility():
    blood_type = ["O-", "O+", "B-", "B+", "A-", "A+", "AB-", "AB+"]
    passed = True
    result_matrix = [[0 for _ in range(8)] for _ in range(8)]
    total_error = []
    for i, recipient in enumerate(blood_type):
        for j, donor in enumerate(blood_type):
            params = GSWParams(q=134217728,n=16, m = 16)
            gsw = GSWScheme(params)
            A, v = gsw.generate_keys() # Pretend Bob can only access A
            a = Alice(convert_type_to_bitstring(recipient), gsw)
            b = Bob(convert_type_to_bitstring(donor), gsw)
            b.A_ctexts = a.share_ctext()
            z, error = a.decrypt(b.eval_circuit())
            total_error.append(error)
            print(z)
            result_matrix[i][7-j] = z
            if z != blood_compatibility_boolean(recipient, donor):
                passed = False
                print(f"[*] Test blood type compatibility failed for, recipient: {recipient}, donor: {donor}")
    print()
    print(np.array(result_matrix))
    if passed:
        print("[*] Passed blood compatibility test")
        print("Average error: ", sum(total_error)/len(total_error))

test_blood_type_compatibility()