from scipy.optimize import minimize
import numpy as np
import time, random

_RATE = 0.18

def computeTarget(todo_btc, todo_eth, todo_usdt, argsr, argstt, argsa):
    #print(repr(argsr),repr(argstt))
    A, B ,C, D, E, F, G ,H, I, J, K, L, M = argsr
    aa, bb, cc, dd, ee, ff, gg, hh, ii, jj, kk, ll, mm = argstt
    Aa, Bb, Cc, Dd, Ee, Ff, Gg, Hh, Ii, Jj, Kk, Ll, Mm = argsa
    print(Aa, Bb, Cc, Dd, Ee, Ff, Gg, Hh, Ii, Jj, Kk, Ll, Mm)

    if Kk >= _RATE:
        _sub1 = (365* G / _RATE )- gg
        btc_bsc = todo_btc * (_sub1 ) / (_sub1  +  (365* K / _RATE )- kk )
    else:
        btc_bsc = todo_btc * 1
    

    _sub2 = (365* H / _RATE )- hh
    _sub22 = _sub2
    for (Z, zz, Zz)  in ( (I, ii, Ii), (J, jj, Jj) ):
        _sub22 += ((365* Z / _RATE )- zz) if  Zz >= _RATE else 0


    eth_bsc = todo_eth * _sub2 / _sub22
    #eth_bsc = todo_eth * ((365* H / _RATE )- hh) / ((365* H / _RATE)- hh + (365* I / _RATE )- ii  + (365* J / _RATE )- jj )
 
    _sub3 = 0
    for (Z, zz, Zz) in ((C, cc, Cc), (D, dd, Dd), (F, ff, Ff), (G, gg, Gg), (H, hh, Hh)):
        _sub3  += ((365* Z / _RATE )- zz) if  Zz >= _RATE else 0

    _sub33 = _sub3
    for (Z, zz, Zz) in ((J, jj, Jj), (M, mm, Mm)):
        _sub33 += ((365* Z / _RATE )- zz) if  Zz >= _RATE else 0
    #_sub = (365* C / _RATE )- cc + (365* D / _RATE)- dd +(365* F / _RATE )- ff  +(365* G / _RATE )- gg +(365* H / _RATE )- hh
    #usdt_bsc= todo_usdt *( _sub ) / ( _sub +(365* J / _RATE )- jj +(365* M / _RATE )- mm)
    usdt_bsc =todo_usdt * _sub3 / _sub33

    return btc_bsc, eth_bsc, usdt_bsc


def totalReward(argsp, argsr, argst, argsa):
    bnb, cake, btcb, eth, busd, usdt = argsp
    #bnb, busd, usdt, cake, btcb, eth = argsp
    A,  B , C,  D,  E,  F,  G ,  H , *_ = argsr
    a,  b,  c,  d,  e,  f,  g ,  h , *_ = argst
    Aa, Bb, Cc, Dd, Ee, Ff, Gg , Hh, *_ = argsa

    func = lambda x: (-1)*(
            (((bnb*x[0] + busd*x[3])/(bnb*x[0] + busd*x[3] + a )*A) if Aa>=_RATE and x[0]>0 and x[3]>0 else (0)) + \
            (((bnb*x[1] + busd*x[2])/(bnb*x[1] + busd*x[2] + b )*B) if Bb>=_RATE and x[1]>0 and x[2]>0 else (0)) + \
            (((bnb*x[5] + usdt*x[11])/(bnb*x[5] + usdt*x[11] + c )*C) if Cc>=_RATE and x[5]>0 and x[11]>0 else (0)) + \
            (((bnb*x[6] + usdt*x[12])/(bnb*x[6] + usdt*x[12] + d )*D) if Dd>=_RATE and x[6]>0 and x[12]>0 else (0)) + \
            (((cake*x[15] + busd*x[4])/(cake*x[15] + busd*x[4] + e )*E)   if Ee >=_RATE and x[15]>0 and x[4]>0 else (0)) + \
            (((cake*x[14] + usdt*x[13])/(cake*x[14] + usdt*x[13] + f )*F) if Ff >=_RATE and x[14]>0 and x[13]>0 else (0)) + \
            (((btcb*x[7] + usdt*x[10])/(btcb*x[7] + usdt*x[10] + g )*G)   if Gg >=_RATE and x[7]>0 and x[10]>0 else (0)) + \
            (((eth*x[8] + usdt*x[9])/(eth*x[8] + usdt*x[9] + h )*H)   if Hh >=_RATE and x[8]>0 and x[9]>0 else (0))
          )
    return func

LOW_BOUND = 0

def conSamllRe(argsq, argsp):
    bnb_q, busd_q, btcb_q, eth_q, usdt_q, cake_q = argsq
    bnb, cake, btcb, eth, busd, usdt = argsp
    #bnb, busd, usdt, cake, btcb, eth = argsp
    # 约束条件 分为eq 和ineq
    # eq表示 函数结果等于0 ； ineq 表示 表达式大于等于0  
    cons = ({'type': 'ineq', 'fun': lambda x: x[0] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: bnb_q - x[0]},
            {'type': 'ineq', 'fun': lambda x: x[1]},
            {'type': 'ineq', 'fun': lambda x: bnb_q - x[1]},
            {'type': 'ineq', 'fun': lambda x: x[2] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: busd_q - x[2]},
            {'type': 'ineq', 'fun': lambda x: x[3] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: busd_q - x[3]},
            {'type': 'ineq', 'fun': lambda x: x[4] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: busd_q - x[4]},
            {'type': 'ineq', 'fun': lambda x: x[5] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: bnb_q - x[5]},
            {'type': 'ineq', 'fun': lambda x: x[6] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: bnb_q - x[6]},
            {'type': 'ineq', 'fun': lambda x: x[7] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: btcb_q - x[7]},
            {'type': 'ineq', 'fun': lambda x: x[8] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: eth_q - x[8]},
            {'type': 'ineq', 'fun': lambda x: x[9] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: usdt_q - x[9]},
            {'type': 'ineq', 'fun': lambda x: x[10] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: usdt_q - x[10]},
            {'type': 'ineq', 'fun': lambda x: x[11] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: usdt_q - x[11]},
            {'type': 'ineq', 'fun': lambda x: x[12] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: usdt_q - x[12]},
            {'type': 'ineq', 'fun': lambda x: x[13] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: usdt_q - x[13]},
            {'type': 'ineq', 'fun': lambda x: x[14] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: cake_q - x[14]}, 
            {'type': 'ineq', 'fun': lambda x: x[15] - LOW_BOUND},
            {'type': 'ineq', 'fun': lambda x: cake_q - x[15]}, 
            {'type': 'ineq', 'fun': lambda x: usdt_q - x[9] - x[10] - x[11] - x[12] - x[13]},  #x[9]???????
            {'type': 'ineq', 'fun': lambda x: busd_q - x[2] - x[3] - x[4]}, 
            {'type': 'ineq', 'fun': lambda x: bnb_q - x[0] - x[1]- x[5] - x[6]}, 
            {'type': 'ineq', 'fun': lambda x: cake_q - x[14] - x[15]}, 
            {'type': 'ineq', 'fun': lambda x: btcb_q - x[7]}, 
            {'type': 'ineq', 'fun': lambda x: eth_q - x[8]}, 
            {'type': 'eq',   'fun': lambda x: bnb*x[0] - busd*x[3]}, 
            {'type': 'eq',   'fun': lambda x: bnb*x[1] - busd*x[2]},
            {'type': 'eq',   'fun': lambda x: bnb*x[5] - usdt*x[11]},
            {'type': 'eq',   'fun': lambda x: bnb*x[6] - usdt*x[12]},
            {'type': 'eq',   'fun': lambda x: cake*x[15] - busd*x[4]},
            {'type': 'eq',   'fun': lambda x: cake*x[14] - usdt*x[13]},
            {'type': 'eq',   'fun': lambda x: btcb*x[7] - usdt*x[10]},
            {'type': 'eq',   'fun': lambda x: eth*x[8] - usdt*x[9]},
            )
    return cons

def getRandX0(bounds):
    x = []
    for b in bounds:
        x.append( random.random() * b )
    return x

def minimizeWrapper(x0, argsq, argsp, argsr, argst, argsa, boundary):
    res_list = []
    cons = conSamllRe(argsq, argsp)
    for m in (#'Nelder-Mead',
              #'Powell', 
              #'CG', 
              #'BFGS',
              #'L-BFGS-B',
              #'TNC',
              'SLSQP',
            ) :
        ss = time.time()
        result = minimize( totalReward(argsp, argsr, argst, argsa), 
                           x0, 
                           method=m, 
                           bounds=tuple([(LOW_BOUND, b)for b in boundary]),
                           constraints=cons)
        res_tup = (result.fun, m, x0, result.x)
        if result.success:
            res_list.append(res_tup)
        else:
            print('not success:', res_tup)
    res_list.sort(key=lambda k:k[0], reverse=True)
    if len(res_list) > 1 :
        if res_list[-1][0]-res_list[0][0] > 0.000000000001:
            print('max diff val of all methods:', res_list[-1][0]-res_list[0][0] )
    return res_list[-1] if len(res_list) > 0 else None

LOOP_COUNT = 100
def doCompute(argsq, argsp, argsr, argst, argsa):
    ss = time.time()

    bnb_q, busd_q, btcb_q, eth_q, usdt_q, cake_q = argsq
    bnb, cake, btcb, eth, busd, usdt = argsp
    A, B ,C, D, E, F, G, H, *_= argsr
    aa, bb, cc, dd, ee, ff, gg, hh, *_ = argst
    Aa, Bb, Cc, Dd, Ee, Ff, Gg, Hh, *_ = argsa
    
    boundary = (bnb_q, bnb_q, busd_q, busd_q, busd_q, bnb_q, bnb_q, btcb_q,
                eth_q, usdt_q, usdt_q, usdt_q, usdt_q, usdt_q, cake_q, cake_q)
    res_list = []

    for i in range(LOOP_COUNT):
        res = minimizeWrapper(getRandX0(boundary), argsq, argsp, argsr, argst, argsa, boundary)
        if res:
            res_list.append(res)

    print("len of res list:",len(res_list))
    print("Time cost: %0.3f"%(time.time()-ss))

    res_list.sort(key=lambda k:k[0], reverse=True)
    for res in res_list:
        #print('Time cost: %0.6f'%(time.time()-ss))
        #print('Method:',m)
        print('Value: %0.6f'%res[0])
        print('Method:', res[1])
        print('Init  X:',res[2])
        print('Final X:',list(map(float, res[3])))
        print('================================================================')
    print('Boundary:', boundary)
    print("argsa:", Aa, Bb, Cc, Dd, Ee, Ff, Gg, Hh)
    if len(res_list)==0 :
        return []

    X = list(map(float, res_list[-1][3]))
    Val = res_list[-1][0]
    print("""
    total info: btc:{} eth:{}, cake:{} bnb:{} usdt:{}
    success: {}
    message: {}
    value: {}
    price btc:{} eth:{} cake:{} bnb:{} usdt:{}
    biswap tvl:{} apr:{} btc:{} usdt:{}
    biswap tvl:{} apr:{} eth:{} usdt:{}
    biswap tvl:{} apr:{} bnb:{} usdt:{}
    pancake tvl:{} apr:{} bnb:{} usdt:{}
    pancake tvl:{} apr:{} cake:{} usdt:{}
    total invest btc:{} eth:{} cake:{} bnb:{} usdt:{}
    diff btc:{} eth:{} cake:{} bnb:{} usdt:{}
    """.format(
        btcb_q , eth_q, cake_q, bnb_q, usdt_q,
        True,
        "",
        Val,
        btcb, eth, cake, bnb, usdt,
        gg, Gg, X[7], X[10],
        hh, Hh, X[8], X[9],
        cc, Cc, X[5], X[11],
        dd, Dd, X[6], X[12],
        ff, Ff, X[14], X[13],
        X[7], X[8], X[14], X[5] + X[6],
        X[9] + X[10] + X[11] + X[12] + X[13],
        btcb_q - X[7],
        eth_q  - X[8],
        cake_q - X[14],
        bnb_q- X[5] - X[6],
        usdt_q - (X[9] + X[10] + X[11] + X[12] + X[13])
        )
    )

    return X

if __name__ == '__main__':
    pass


