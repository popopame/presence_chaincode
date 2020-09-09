package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


//Definition of a SmartContract that contain the Contract API , we will build our function on this struct
type SmartContract struct {
	contractapi.Contract
}

//Definition of a Identity Struct , that contain all the information that we want to validate in the World State. The ID is used as the key 
type Identity struct {
	ID             string `json:"ID"`
	Surname		   string `json:"Surname"`
	Name           string `json:"name"`
	NameTag        string `json:"nametag"`
	TwitterID      string `json:"twitterid"`
	DiscordID	   string `json:"discordid"`
	InstagramID    string `json:"instagramid"`
	TwitchID       string `json:"twitchid"`
}

//Init Ledger marshall some stub ID and put them in the world state. The ID is used as the key 
func (sc *SmartContract) InitLedger( ctx contractapi.TransactionContextInterface) error {
	StubID := []Identity{
		{ID: "JD2020230001",Surname: "John",Name: "Doe",NameTag:"SuperUser",TwitterID:"TwitterID39",DiscordID:"DiscordID92",InstagramID:"InstagramID02",TwitchID:"SKDIAS"},
		{ID: "OD2020230001",Surname: "Oliver",Name: "Tree",NameTag:"Pwned",TwitterID:"TwitterIDF39",DiscordID:"DiscordID03",InstagramID:"InstagramID32",TwitchID:"SKDSSAS"},
	}
	for _,Identity := range StubID {
		IdentityJson, err := json.Marshal(Identity)

		if err != nil {
			return err
		}

		ctx.GetStub().PutState(Identity.ID,IdentityJson)

		if err != nil {
			return fmt.Errorf("Failed to set inital value on worldState. Returned error: %s", err)
		}
	}

	return nil 
}

//createID function is used to add an ID to the World state , it take all the elements of and ID , Marshall them and put them in the world state
func (sc *SmartContract) createID(ctx contractapi.TransactionContextInterface, ID string , surname string , name string , nameTag string , twitterID string , discordID string , instagramID string , twitchID string) error {
	
	Identity :=Identity{
		ID: ID ,
		Surname: surname,
		Name: name,
		NameTag: nameTag,
		TwitterID: twitterID,
		DiscordID: discordID,
		InstagramID: instagramID,
		TwitchID: twitchID,		
	}
	existingID , err :=ctx.GetStub().GetState(ID)

	if err != nil {
		return fmt.Errorf("Failed to interact with world state , error returned: %s" , err)
		
	}

	if existingID != nil {
		return fmt.Errorf("Cannot create Identity with Given ID , one already present in the world state with value: %s", existingID)
	}


	IdentityJson, err := json.Marshal(Identity)

	if err != nil {
		return err
	}

	ctx.GetStub().PutState(Identity.ID,IdentityJson)

	if err != nil {
		return fmt.Errorf("Failed to set ID in World State , error returned: %s", err)
	}
	return nil

}

//getID function will query the Identity with a given key
func (sc *SmartContract) getID(ctx contractapi.TransactionContextInterface, ID string) (*Identity,error) {

	identityAsByte , err := ctx.GetStub().GetState(ID)
	if err != nil {
		return nil, fmt.Errorf("Can't interact with world state , returned error : %s" ,err)
	}

	if identityAsByte == nil {
		return nil , fmt.Errorf("Can't retrieve specified ID , value does not exist")
	}

	identity := new(Identity)

	_ = json.Unmarshal(identityAsByte,identity)

	return identity, nil
}

// changeValue func will change a value in a Identity , it return the error if one occurs
func (s *SmartContract) changeNameTag(ctx contractapi.TransactionContextInterface , ID string, newValue string) error {
	
	identity , err := s.getID(ctx,ID)

	if err != nil {
		return fmt.Errorf("Can't interact with world state , returned error : %s" ,err)
	}

	if identity == nil {
		return fmt.Errorf("Can't find any identity with ID %s in the world state. erro: %s",ID,err)
	}

	identity.NameTag = newValue
	identityAsByte, _ := json.Marshal(identity)

	return ctx.GetStub().PutState(ID,identityAsByte)

}


func main() {
	chaincode , err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Can't creating chaincode: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode , returned error: %s",err)
	}

}
