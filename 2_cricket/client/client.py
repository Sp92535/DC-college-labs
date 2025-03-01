import sys
import os

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

import grpc
import cricket_pb2
import cricket_pb2_grpc

def run():
    # Connect to the gRPC server
    channel = grpc.insecure_channel('localhost:6969')
    stub = cricket_pb2_grpc.CricketStub(channel)

    # Call GetTopScorers
    response = stub.getTopScorers(cricket_pb2.Empty())
    print(f"Top Scorer: {response.name}, Average: {response.average}")

    # Call GetCenturions
    response = stub.getCenturions(cricket_pb2.Empty())
    print(f"Most Centuries: {response.name}, Centuries: {response.centuries}")

    # Call GetPlayerStats
    player_name = "Virat Kohli"
    response = stub.getPlayerStats(cricket_pb2.PlayerRequest(name=player_name))
    print(f"Stats for {response.name} - Average: {response.average}, Centuries: {response.centuries}")

    # Call UpdatePlayerScore
    stub.updatePlayerScore(cricket_pb2.UpdateScoreRequest(name="Virat Kohli", runs=100))
    print("Updated player score successfully!")

if __name__ == "__main__":
    run()
