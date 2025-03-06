from flask import Flask, jsonify

app = Flask(__name__)


def solve_n_queens(n):
    def DFS(queens, xy_dif, xy_sum):
        p = len(queens)
        if p==n:
            result.append(queens)
            return None
        for q in range(n):
            if q not in queens and p-q not in xy_dif and p+q not in xy_sum: 
                DFS(queens+[q], xy_dif+[p-q], xy_sum+[p+q])  
    result = []
    DFS([],[],[])
    return [ ["."*i + "Q" + "."*(n-i-1) for i in sol] for sol in result]

@app.route("/<int:n>")
def n_queens_route(n):
    solution = solve_n_queens(n)
    return jsonify({"n": n, "solution": solution})


if __name__ == "__main__":
    app.run()
