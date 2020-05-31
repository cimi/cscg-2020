import gmpy2
gmpy2.get_context().precision = 4096

from binascii import unhexlify
from functools import reduce
from gmpy2 import root
from Crypto.Util.number import long_to_bytes

# Hastad's Broadcast Attack
# https://id0-rsa.pub/problem/11/

# Resources
# https://en.wikipedia.org/wiki/Coppersmith%27s_Attack
# https://github.com/sigh/Python-Math/blob/master/ntheory.py

EXPONENT = 3

CIPHERTEXT_1 = 3999545484320691620582760666106855727053549021662410570083429799334896462058097237449452993493720397790227435476345796746350169898032571754431738796344192821893497314910675156060408828511224220581582267651003911249219982138536071681121746144489861384682069580518366312319281158322907487188395349879852550922320727712516788080905540183885824808830769333571423141968760237964225240345978930859865816046424226809982967625093916471686949351836460279672029156397296634161792608413714942060302950192875262254161154196090187563688426890555569975685998994856798884592116345112968858442266655851601596662913782292282171174885
CIPHERTEXT_2 = 7156090217741040585758955899433965707162947606350521948050112381514262664247963697650055668324095568121356193295269338497644168513453950802075729741157428606617001908718212348868412342224351012838448314953813036299391241983248160741119053639242636496528707303681650997650419095909359735261506378554601448197330047261478549324349224272907044375254024488417128064991560328424530705840832289740420282298553780466036967138660308477595702475699772675652723918837801775022118361119700350026576279867546392616677468749480023097012345473460622347587495191385237437474584054083447681853670339780383259673339144195425181149815
CIPHERTEXT_3 = 9343715678106945233699669787842699250821452729365496523062308278114178149719235923445953522128410659220617418971359137459068077630717894445019972202645078435758918557351185577871693207368250243507266991929090173200996910881754217374691865096976051997491208921880703490275111904577396998775470664002942232492755888378994040358902392803421017545356248082413409915177589953816030273082416979477368273328755386893089597798104163894528521114946660635364704437632205696975201216810929650384600357902888251066301913255240181601332549134854827134537709002733583099558377965114809251454424800517166814936432579406541946707525

MODULUS_1 = 13368388890946686131727968139222647635627171995393331225756908262294343216259723081458150905003600322756476137516299938365001972798137046621672975975457070465770187049834603521354462199081700902700733323087201964703391196426066952717511505120664658507099276380167252844734836468387820963170177521935571096868999638202790415914397116993003197932961837711222659120426461631108658146330240545816025557486272830688061978425683447522103977339616076727857816034089500594682018085999092378789197039633371210351470521621878994691517983319668541047042031499811379908242466040735576388227260217406960791406632454767045448789863
MODULUS_2 = 14463626602170229427356167809091927075048214837573339781774138582390190263460223568524802570585480435667949138330700031482283411314199309061664373861923286634420548935259474128834717819970239283387732315996647938605905994532994027238099470750924616969478147212529380894358056424265545387574975098446117146942068553320197224781384410276446833888437566192029289304444125818681142352673878184276408904704077528699342956063922184456200815444422094356292649411256904543442078043661428831462400371961179888731725665328211651272084619341652653674440701885337593085045665899694222470709757866143325441669946933338126683188131
MODULUS_3 = 18343717802716940630601940481023526351437486074120550591161058762020703345710367605696446334690825248791560509373517279950125583944976720622084902078751153032339436688975255343139180338681178127688700797071320999813387670292350135483485169318320316776584245519471849328634745109353977968597175721348420576770527793000136160877295577014905354451575371196006765377541964045640054268423795610241643005381587433138330817893094851452345761462684724873155990606241842996499888181450611803912139827073505685135888393196549213527418583778495818537291115829823762105372358484486446314835437285354604977091862400207219042791731


def chinese_remainder_theorem(items):
    # Determine N, the product of all n_i
    N = 1
    for a, n in items:
        N *= n

    # Find the solution (mod N)
    result = 0
    for a, n in items:
        m = N // n
        r, s, d = extended_gcd(n, m)
        if d != 1:
            raise "Input not pairwise co-prime"
        result += a * s * m

    # Make sure we return the canonical solution.
    return result % N


def extended_gcd(a, b):
    x, y = 0, 1
    lastx, lasty = 1, 0

    while b:
        a, (q, b) = b, divmod(a, b)
        x, lastx = lastx - q * x, x
        y, lasty = lasty - q * y, y

    return (lastx, lasty, a)


def mul_inv(a, b):
    b0 = b
    x0, x1 = 0, 1
    if b == 1:
        return 1
    while a > 1:
        q = a // b
        a, b = b, a % b
        x0, x1 = x1 - q * x0, x0
    if x1 < 0:
        x1 += b0
    return x1


def get_value(filename):
    with open(filename) as f:
        value = f.readline()
    return int(value, 16)

if __name__ == '__main__':

    # C1 = get_value(CIPHERTEXT_1)
    # C2 = get_value(CIPHERTEXT_2)
    # C3 = get_value(CIPHERTEXT_3)

    C1 = CIPHERTEXT_1
    C2 = CIPHERTEXT_2
    C3 = CIPHERTEXT_3
    ciphertexts = [C1, C2, C3]

    # N1 = get_value(MODULUS_1)
    # N2 = get_value(MODULUS_2)
    # N3 = get_value(MODULUS_3)

    N1 = MODULUS_1
    N2 = MODULUS_2
    N3 = MODULUS_3

    modulus = [N1, N2, N3]

    C = chinese_remainder_theorem([(C1, N1), (C2, N2), (C3, N3)])
    M = int(root(C, 3))

    # M = hex(M)[2:]
    # print(unhexlify(M).decode('utf-8'))
    print(long_to_bytes(M))
