import sys
import os
import time
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
import grpc
import cricket_pb2
import cricket_pb2_grpc


def measure_time(func, *args):
    start_time = time.time()
    result = func(*args)
    elapsed_time = (time.time() - start_time) * 1000  # Convert to milliseconds
    return result, elapsed_time

def run():
    # Connect to the gRPC server
    channel = grpc.insecure_channel('localhost:6969')
    stub = cricket_pb2_grpc.CricketStub(channel)

    # Call GetTopScorers
    response, time_taken = measure_time(stub.getTopScorers, cricket_pb2.Empty())
    print(f"Top Scorer: {response.name}, Average: {response.average} (Time taken: {time_taken:.2f} ms)")

    # Call GetCenturions
    response, time_taken = measure_time(stub.getCenturions, cricket_pb2.Empty())
    print(f"Most Centuries: {response.name}, Centuries: {response.centuries} (Time taken: {time_taken:.2f} ms)")

    # Call GetPlayerStats
    player_name = "Virat Kohli"
    response, time_taken = measure_time(stub.getPlayerStats, cricket_pb2.PlayerRequest(name=player_name))
    print(f"Stats for {response.name} - Average: {response.average}, Centuries: {response.centuries} (Time taken: {time_taken:.2f} ms)")

    # Call UpdatePlayerScore
    _, time_taken = measure_time(stub.updatePlayerScore, cricket_pb2.UpdateScoreRequest(name="Virat Kohli", runs=100))
    print(f"Updated player score successfully! (Time taken: {time_taken:.2f} ms)")

if __name__ == "__main__":
    run()
